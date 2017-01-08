package api

import (
	"fmt"
	"net/http"

	"gopkg.in/redis.v5"

	"github.com/gorilla/websocket"
)

var pub *redis.IntCmd

// type Event struct {
// 	Type    string
// 	JobID   string
// 	Message string
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (E *Engine) pubsub(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()

		if string(message) == "subscribe" {

			E.NewSubscriber("events", func(channel string, payload string) {
				// var e Event

				// err = json.Unmarshal([]byte(payload), &e)

				// if err != nil {
				// 	log.Printf("Unmarshal error: %v", err)
				// }

				conn.WriteMessage(mt, []byte(payload))
			})
		}

		// if err != nil {
		// 	fmt.Println("read:", err)
		// 	break
		// }
		//fmt.Printf("recv: %s\n", message)
		//err = conn.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("Error:", err)
			break
		}
	}
}
