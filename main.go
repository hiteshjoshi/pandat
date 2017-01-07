package main

import (
	"github.com/hiteshjoshi/pandat/api"
	"github.com/hiteshjoshi/pandat/clock"
)

func main() {

	c := clock.New()
	go c.Boot()
	api.Start("9090", c)
}
