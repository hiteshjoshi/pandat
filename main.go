package main

import (
	"fmt"

	"github.com/hiteshjoshi/pandat/api"
	"github.com/hiteshjoshi/pandat/clock"
)

var (
	Version string
	Build   string
)

func main() {

	fmt.Println("Version: ", Version)
	fmt.Println("Build: ", Build)

	c := clock.New()
	go c.Boot() //this makes sure that we have persistent connection to clock server
	api.Start("9090", c)
}
