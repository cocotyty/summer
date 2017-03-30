package main

import (
	"github.com/cocotyty/summer"
	"time"
)

func main() {
	summer.Start()
	waitToExit()
}
func waitToExit() {
	for {
		time.Sleep(time.Hour)
	}
}
