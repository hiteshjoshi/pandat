package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hiteshjoshi/pandat/clock"
	validator "gopkg.in/validator.v2"
)

type EventController struct {
	Clock *clock.Clock
}

func (E *EventController) Get(w http.ResponseWriter, r *http.Request) {

	resp := Response{
		Message: "Events attached.",
		Data:    map[string]string{},
		Error:   false,
	}
	resp.Send(http.StatusOK, w)
}

func (E *EventController) Add(w http.ResponseWriter, r *http.Request) {

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

	id := E.Clock.Add(schedule.Interval, schedule.URL)

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

func (E *EventController) Remove(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventID := vars["eventID"]

	err := E.Clock.Remove(eventID)

	if err != nil {
		resp := Response{
			Message: "Server error",
			Error:   true,
		}
		resp.Send(http.StatusBadRequest, w)
		return
	}

	resp := Response{
		Message: "Event removed.",
		Error:   false,
	}
	resp.Send(http.StatusOK, w)
}
