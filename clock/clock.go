package clock

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"strconv"

	"os"

	"gopkg.in/redis.v5"
)

const (
	USERAGENT       = "Pandat API"
	TIMEOUT         = 30
	CONNECT_TIMEOUT = 5
)

type fn func() string
type Clock struct {
	ID    uint16
	Cron  *Cron
	Boot  fn
	Redis *redis.Client
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS"),
		PoolSize: 100,
	})
}

func (c *Clock) forever() {
	for {
		_, err := c.Redis.Ping().Result()

		//fmt.Println("PING TO REDIS, GOT : ", pong, "and ERRORS : ", err)
		if err != nil {
			c.Redis = NewRedisClient()
		}
		time.Sleep(time.Second * 5)
	}
}
func New() *Clock {

	c := NewCron()

	redisClient := NewRedisClient()

	pong, err := redisClient.Ping().Result()

	fmt.Println("PING TO REDIS, GOT : ", pong, "and ERRORS : ", err)

	//CLOCK SERVER ENTRIES
	clockID, _ := redisClient.Incr("clocks").Result()

	fmt.Println("CLock ID ", clockID)

	clock := Clock{
		Cron:  c,
		Redis: redisClient,
		ID:    uint16(clockID),
	}

	//CRON Entries
	entries, _ := redisClient.Get("entries").Result()
	id, _ := strconv.Atoi(entries)
	c.NextID = EntryID(id)

	//To make it run forever for slaves of clock
	clock.Boot = func() string {
		done := make(chan bool)
		go clock.forever()
		<-done
		return ""
	}

	c.Start() //start cron.v2

	return &clock
}

func (c *Clock) Stop() error {
	_, err := c.Redis.IncrBy("clocks", -1).Result()
	c.Redis.Close()
	return err
}

//Add : Add new job
func (c *Clock) Add(interval string, url string) string {

	id, _ := c.Cron.AddJob(interval, Request{
		URL:     url,
		Redis:   c.Redis,
		ClockID: c.ID,
	})

	ID := fmt.Sprint(id)

	command := fmt.Sprint("cron:", c.ID, "::entry::", ID)

	//save entryID
	c.Redis.HSet(command, "id", ID)
	//save interval
	c.Redis.HSet(command, "interval", interval)
	//save url
	c.Redis.HSet(command, "url", url)
	//empty last run
	c.Redis.HSet(command, "last_run", "")
	//when did we create this one?
	c.Redis.HSet(command, "created", time.Now().String())

	c.Redis.Set("entries", ID, 0)

	return ID

}

//Remove a job
func (c *Clock) Remove(id string) error {

	ID, _ := strconv.Atoi(id)

	c.Cron.Remove(EntryID(ID))

	//remvove from redis
	e := c.Redis.HDel("id:" + id)

	return e.Err()

}

func (c *Clock) GetAll() ([]map[EntryID]string, error) {

	entries, err := c.
		Redis.
		Keys("cron:" + fmt.Sprint(c.ID, "*")).Result()

	if err == nil {
		var clockEntries []map[EntryID]string

		for _, key := range entries {
			fmt.Println(key)
			// entryDetail, _ := c.
			// 	Redis.
			// 	HGetAll(key).
			// 	Result()

			// clockEntries = append(clockEntries, interface{
			// 	"id":entryDetail["id"],
			// })
		}
		return clockEntries, nil
	}

	return nil, err
}

func (c *Clock) GetAllString() ([]map[string]string, error) {
	entries, err := c.
		Redis.
		Keys("cron:" + fmt.Sprint(c.ID, "*")).Result()

	if err == nil {
		var clockEntries []map[string]string

		for _, key := range entries {

			entryDetail, _ := c.
				Redis.
				HGetAll(key).
				Result()

			clockEntries = append(clockEntries, entryDetail)
		}
		return clockEntries, nil
	}

	return nil, err
}

//This function controls what to run on cron execution
type Request struct {
	URL     string
	Redis   *redis.Client
	ClockID uint16
	EntryID string
}

func (d Request) Exec(ch chan<- string) {
	start := time.Now()

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	var netClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest("POST", d.URL, strings.NewReader("data"))
	resp, err := netClient.Do(req)
	req.Close = true
	resp.Close = true
	defer resp.Body.Close()

	if err != nil {
		//TODO create an error for this event
		log.Print("Unable to create a new http request", err)
	}
	//save status in db
	log.Printf((d.URL)+" %s", resp.Status)

	secs := time.Since(start).Seconds()

	command := fmt.Sprint("cron:", d.ClockID, "::entry::", d.EntryID)

	d.Redis.HSet(command, "last_run", time.Now().String())

	d.Redis.HSet(command, "last_meta", map[string]string{
		"time": fmt.Sprint(secs),
	})

	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, d.URL)

}
func (d Request) Run() {
	ch := make(chan string)
	go d.Exec(ch)
}
