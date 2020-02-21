package main

import "os"

func main() {
	args := os.Args
	if args[1] == "cloudwatch" {
		cloudwatchAgentStart(args[2], args[3], args[4], args[5])
	} else {
		println("Command does not exist")
	}
}
