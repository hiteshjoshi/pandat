package clock

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"gopkg.in/redis.v5"
)

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
