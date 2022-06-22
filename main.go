package main

import (
	"exercism-mentoring-request-notifier/collector"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/distributor"
	"exercism-mentoring-request-notifier/files"
	"exercism-mentoring-request-notifier/logging"
	"exercism-mentoring-request-notifier/request"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	defaultConfigPath = "config.json"
	defaultCachePath  = "cache.json"
)

func main() {
	configPath := flag.String("config", defaultConfigPath, "Defines the path to the config")
	cacheFilePath := flag.String("cache", defaultCachePath, "Defines the path to the cache.json")
	printVersionInfo := flag.Bool("v", false, "Defines if the version of the current binary should be printed and then exited")
	flag.Parse()

	if *printVersionInfo {
		printVersionAndExit()
	}

	var cfg = &config.Config{}

	err := files.New(os.ReadFile).JSONToStruct(*configPath, cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	var chMentoringRequests = make(chan map[string][]request.MentoringRequest, 5)

	dist, err := distributor.New(cfg, chMentoringRequests, *cacheFilePath)
	if err != nil {
		log.Fatalf("failed to setup distributor: %v", err)
	}
	err = dist.StartupCheck()
	if err != nil {
		log.Fatalf("startup-check of distributor failed: %s", err.Error())
	}
	go dist.Run()

	col, err := collector.New(cfg, chMentoringRequests, logging.SetupLogging)
	if err != nil {
		log.Fatalf("failed to setup collector: %v", err)
	}

	col.Run()
}
