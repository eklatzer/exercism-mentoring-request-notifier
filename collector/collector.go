package collector

import (
	"exercism-mentoring-request-notifier/collector/client"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/request"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	logFile = "collector_log.json"
)

type Collector struct {
	config       *config.Config
	chanRequests chan map[string][]request.MentoringRequest
	log          *logrus.Logger
}

func New(cfg *config.Config, chRequests chan map[string][]request.MentoringRequest, setupLogging func(logger *logrus.Logger, level, path string) error) (*Collector, error) {
	var c = &Collector{
		config:       cfg,
		chanRequests: chRequests,
		log:          &logrus.Logger{},
	}

	err := setupLogging(c.log, cfg.LogLevel, logFile)
	return c, err
}

func (d *Collector) Run() {
	var httpClient = client.ExercismHTTPClient{Token: d.config.ExercismToken, Client: &http.Client{}}
	for {
		time.Sleep(time.Duration(d.config.Interval) * time.Second)
		results, err := httpClient.GetMentoringRequestsForAllTracks(d.config.TrackConfig)
		if err != nil {
			d.log.Warn(err.Error())
			continue
		}
		d.chanRequests <- results
	}
}
