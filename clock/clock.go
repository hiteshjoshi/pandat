package clock

import (
	"fmt"
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

type fn func()
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
	clock.Boot = func() {
		done := make(chan bool)
		go clock.forever()
		<-done
	}

	c.Start() //start cron.v2

	return &clock
}

func (c *Clock) Stop() error {
	_, err := c.Redis.IncrBy("clocks", -1).Result()
	c.Redis.Close()
	return err
}
