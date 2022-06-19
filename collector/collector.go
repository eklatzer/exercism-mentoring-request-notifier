package collector

import (
	"encoding/json"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/logging"
	"exercism-mentoring-request-notifier/mentoring_request"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	exercismAPIBasePath      = "https://exercism.org/api/v2"
	getMentoringRequestsPath = "/mentoring/requests?track_slug=%s"
	logFile                  = "collector_log.json"
)

type Collector struct {
	config       *config.Config
	chanRequests chan map[string][]mentoring_request.MentoringRequest
	log          *logrus.Logger
}

func New(cfg *config.Config, chRequests chan map[string][]mentoring_request.MentoringRequest) (*Collector, error) {
	var c = &Collector{
		config:       cfg,
		chanRequests: chRequests,
		log:          &logrus.Logger{},
	}

	err := logging.SetupLogging(c.log, cfg.LogLevel, logFile)
	return c, err
}

func (d *Collector) Run() {
	var httpClient = ExercismHttpClient{
		Client: &http.Client{},
		Token:  d.config.ExercismToken,
	}
	for true {
		time.Sleep(time.Duration(d.config.Interval) * time.Second)
		var results = map[string][]mentoring_request.MentoringRequest{}
		for trackSlug := range d.config.TrackConfig {
			requests, err := httpClient.getMentoringRequests(trackSlug)
			if err != nil {
				d.log.Error(err)
				continue
			}
			results[trackSlug] = requests.MentoringRequests
		}
		d.chanRequests <- results
	}
}

func (c *ExercismHttpClient) getMentoringRequests(trackSlug string) (*mentoring_request.MentoringRequestsResults, error) {
	//TODO: pagination
	requestURL := fmt.Sprintf("%s%s", exercismAPIBasePath, fmt.Sprintf(getMentoringRequestsPath, trackSlug))
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http-request failed, status-code: %d, response: %s", resp.StatusCode, body)
	}

	var data = &mentoring_request.MentoringRequestsResults{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
