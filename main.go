package main

import (
	"runtime"

	"github.com/hiteshjoshi/pandat/api"
	"github.com/hiteshjoshi/pandat/clock"
)

func main() {
	runtime.LockOSThread()

	c := clock.New()
	api.Start(8000, c)

}
