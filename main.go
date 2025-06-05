package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

var BuildVersion = "dev"
var startTime time.Time // Track process start time

func main() {
	startTime = time.Now() // Set process start time
	conf := flag.String("config", "config.json", "path to config file or a http(s) url")
	version := flag.Bool("version", false, "print version and exit")
	help := flag.Bool("help", false, "print help and exit")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Println(BuildVersion)
		return
	}
	config, err := load(*conf)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	err = startHTTPServer(config)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
