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

	cron "gopkg.in/hiteshjoshi/cron.v2"
	"gopkg.in/redis.v5"
)

const (
	USERAGENT       = "Pandat API"
	TIMEOUT         = 30
	CONNECT_TIMEOUT = 5
)

type fn func() string
type Clock struct {
	Cron  *cron.Cron
	Boot  fn
	Redis *redis.Client
}

func forever() {
	for {
	}
}
func New() *Clock {

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS"),
	})
	pong, err := client.Ping().Result()

	fmt.Println("PING TO REDIS, GOT : ", pong, "and ERRORS : ", err)

	c := cron.New()

	e := client.Get("entries")
	s, _ := e.Bytes()
	id, _ := strconv.Atoi(string(s))
	c.NextID = cron.EntryID(id)

	clock := Clock{
		Cron:  c,
		Redis: client,
	}

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

type Job struct {
	URL string
}

func (d Job) Request(ch chan<- string) {
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
func (d Job) Run() {
	ch := make(chan string)
	go d.Request(ch)

}

//Add : Add new job
func (c *Clock) Add(interval string, url string) string {

	id, _ := c.Cron.AddJob(interval, Job{
		URL: url,
	})

	ID := fmt.Sprint(id)

	//save interval
	c.Redis.HSet("id:"+ID, "interval", interval)

	//save url
	c.Redis.HSet("id:"+ID, "url", url)

	c.Redis.Set("entries", ID, 0)

	return ID

}

//Remove a job
func (c *Clock) Remove(id string) error {

	ID, _ := strconv.Atoi(id)

	c.Cron.Remove(cron.EntryID(ID))

	//remvove from redis
	e := c.Redis.HDel("id:" + id)

	return e.Err()

}
