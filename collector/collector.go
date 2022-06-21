package collector

import (
	"exercism-mentoring-request-notifier/collector/client"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/logging"
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

func New(cfg *config.Config, chRequests chan map[string][]request.MentoringRequest) (*Collector, error) {
	var c = &Collector{
		config:       cfg,
		chanRequests: chRequests,
		log:          &logrus.Logger{},
	}

	err := logging.SetupLogging(c.log, cfg.LogLevel, logFile)
	return c, err
}

func (d *Collector) Run() {
	var httpClient = client.ExercismHTTPClient{
		Client: &http.Client{},
		Token:  d.config.ExercismToken,
	}
	for {
		time.Sleep(time.Duration(d.config.Interval) * time.Second)
		var results = map[string][]request.MentoringRequest{}
		for trackSlug := range d.config.TrackConfig {
			requests, err := httpClient.GetAllMentoringRequests(trackSlug)
			if err != nil {
				d.log.Error(err)
				continue
			}
			results[trackSlug] = requests
		}
		d.chanRequests <- results
	}
}
