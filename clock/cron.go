package clock

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	cron "gopkg.in/robfig/cron.v2"
)

const (
	USERAGENT       = "Pandat API"
	TIMEOUT         = 30
	CONNECT_TIMEOUT = 5
)

type Clock struct {
	Cron *cron.Cron
}

func New() *Clock {

	c := cron.New()
	clock := Clock{
		Cron: c,
	}
	c.Start()
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
		log.Print("Unable to create a new http request", err)
	}

	if err != nil {
		log.Println("error POSTing example.com", err)
	}

	log.Printf("example.com %s", resp.Status)

	secs := time.Since(start).Seconds()

	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s", secs, d.URL)

}
func (d Job) Run() {
	ch := make(chan string)
	go d.Request(ch)

}

//Add : Add new job
func (c *Clock) Add(interval string, url string) string {

	job := Job{
		URL: url,
	}

	id, _ := c.Cron.AddJob(interval, job)

	return strconv.Itoa(int(id))

}
