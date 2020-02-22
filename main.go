package main

import (
	"os"
	"time"
)

func main() {
	args := os.Args
	if args[1] == "cloudwatch" {
		for {
			cloudwatchAgentStart(args[2], args[3], args[4], args[5])
			time.Sleep(10 * time.Second)
		}
	} else {
		println("Command does not exist")
	}
}
