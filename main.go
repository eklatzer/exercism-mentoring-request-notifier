package main

import (
	"exercism-mentoring-request-notifier/collector"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/distributor"
	"exercism-mentoring-request-notifier/mentoring_request"
	"flag"
	log "github.com/sirupsen/logrus"
)

const (
	defaultConfigPath = "config.json"
)

func main() {
	configPath := flag.String("config", defaultConfigPath, "Defines the path to the config")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var chMentoringRequests = make(chan []mentoring_request.MentoringRequest, 5)

	col, err := collector.New(cfg, chMentoringRequests)
	if err != nil {
		log.Fatalf("failed to setup collector: %v", err)
	}

	go col.Run()

	dist, err := distributor.New(cfg, chMentoringRequests)
	if err != nil {
		log.Fatalf("failed to setup distributor: %v", err)
	}
	dist.Run()
}
