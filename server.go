package main

import (
	"log"
	"os"
)

var (
	Sentinel *SentinelStat
	Instance *SentinelInstance
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal("config file is required")
	}
	Sentinel = NewSentinelStat(os.Args[1])

	err := Sentinel.HandleConfiguration()
	if err != nil {
		log.Fatal("failed to load sentinel configuration: ", err.Error())
	}

	Sentinel.Run()
}
