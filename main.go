package main

import (
	"github.com/hiteshjoshi/pandat/api"
	"github.com/hiteshjoshi/pandat/clock"
)

func main() {

	c := clock.New()
	api.Start("9090", c)
}
