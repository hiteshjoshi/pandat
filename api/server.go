package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hiteshjoshi/pandat/clock"

	validator "gopkg.in/validator.v2"
)

func Start(PORT int, c *clock.Clock) *Engine {

	r := NewRouter()
	//set cron

	r.Clock = c

	r.Post("/", r.Index)

	r.Post("/add", r.Add)

	r.Run(":8000")

	return r
}

type Schedular struct {
	Name     string `json:"name,required" validate:"min=2,nonzero"`
	URL      string `json:"url,required" validate:"min=2,nonzero"`
	Interval string `json:"interval,required" validate:"min=2,nonzero"`
}

func (e *Engine) Add(w http.ResponseWriter, r *http.Request) {

	schedule := Schedular{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&schedule)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}
	errs := validator.Validate(schedule)

	if errs != nil {
		resp := Response{
			Message: "Validation error",
			Data:    schedule,
			Error:   true,
		}
		resp.Send(http.StatusBadRequest, w)
		return
	}

	id := e.Clock.Add(schedule.Interval, schedule.URL)

	resp := Response{
		Message: "Event created.",
		Data: map[string]string{
			"id":        id,
			"scheduled": schedule.Interval,
		},
		Error: false,
	}
	resp.Send(http.StatusOK, w)
}

func (e *Engine) Index(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL.Query())

	resp := Response{
		Message: "Yo bro, All good?",

		Error: false,
	}
	resp.Send(http.StatusOK, w)
}
