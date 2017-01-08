package main

import (
	"github.com/hiteshjoshi/pandat/api"
	"github.com/hiteshjoshi/pandat/clock"
)

func main() {

	c := clock.New()
	go c.Boot() //this makes sure that we have persistent connection to clock server
	api.Start("9090", c)
}
