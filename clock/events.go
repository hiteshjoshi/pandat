package clock

import (
	"fmt"
	"strconv"
	"time"
)

type Event struct {
	Type    string
	JobID   string
	Message string
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

	//Publish to events
	c.Publish("events", &Event{
		Type:    "added",
		JobID:   ID,
		Message: "New job added",
	})

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
