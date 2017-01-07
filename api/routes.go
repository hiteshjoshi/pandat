package api

import (
	"net/http"

	"github.com/hiteshjoshi/pandat/clock"
)

func (E *Engine) Routes() {

	eventController := EventController{
		Clock: E.Clock,
	}

	E.Get("/", Index)

	//for realtime analytics
	E.WS("/ws", E.pubsub)

	events := E.Group("/events")
	{
		events.Get("/", eventController.Get)
		events.Post("/", eventController.Add)
		events.Post("/{eventID}", eventController.Remove)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {

	resp := Response{
		Message: "Yo bro, All good?",

		Error: false,
	}
	resp.Send(http.StatusOK, w)
}

func Start(PORT string, c *clock.Clock) {

	E := NewEngine()

	E.Clock = c

	E.Routes()

	//run this http engine
	E.Run(":" + PORT)
}
