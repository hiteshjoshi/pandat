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

func forever() {
	for {
	}
}
func New() *Clock {

	c := NewCron()

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS"),
	})
	pong, err := redisClient.Ping().Result()

	fmt.Println("PING TO REDIS, GOT : ", pong, "and ERRORS : ", err)

	//CLOCK SERVER ENTRIES
	clockID, _ := redisClient.Incr("clock_servers").Result()
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
		go forever()
		<-done
		return ""
	}

	c.Start() //start cron.v2

	return &clock
}

func (c *Clock) Stop() error {
	_, err := c.Redis.IncrBy("clock_servers", -1).Result()
	//fmt.Println(i, err)
	return err
}

//Add : Add new job
func (c *Clock) Add(interval string, url string) string {

	id, _ := c.Cron.AddJob(interval, Request{
		URL: url,
	})

	ID := fmt.Sprint(id)

	//save interval
	c.Redis.HSet("id:"+string(c.ID), "interval", interval)

	//save url
	c.Redis.HSet("id:"+string(c.ID), "url", url)

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

func (c *Clock) GetAll() {
	//c.Redis.HGetAll
}

//This function controls what to run on cron execution
type Request struct {
	URL string
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

	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, d.URL)

}
func (d Request) Run() {
	ch := make(chan string)
	go d.Exec(ch)
}
