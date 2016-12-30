package main

import (
	"fmt"
	"os"
	"os/signal"

	"gopkg.in/robfig/cron.v2"
)

func main() {
	c := cron.New()
	c.AddFunc("@every 2s", func() { fmt.Println("Every hour on the half hour") })

	c.Start()

	s := make(chan os.Signal, 2)

	signal.Notify(s, os.Interrupt)
	go func() {
		for range s {
			//stop cron server
			fmt.Println("\nStopping server!")
			c.Stop()

			os.Exit(0)
		}
	}()

	keep()
}

func keep() {

	done := make(chan bool)
	go (func() {
		for {
		}
	})()
	<-done
}
