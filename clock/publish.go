package clock

import (
	"encoding/json"

	"gopkg.in/redis.v5"
)

func (C *Clock) PublishString(channel, message string) *redis.IntCmd {
	return C.Redis.Publish(channel, message)
}

func (C *Clock) Publish(channel string, message interface{}) *redis.IntCmd {
	// TODO reflect if interface{} type is string, Publish as-is
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	messageString := string(jsonBytes)
	return C.Redis.Publish(channel, messageString)
}
