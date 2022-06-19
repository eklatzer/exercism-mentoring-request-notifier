package distributor

import (
	"encoding/json"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/files"
	"exercism-mentoring-request-notifier/logging"
	"exercism-mentoring-request-notifier/mentoring_request"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"io/ioutil"
	"os"
)

const (
	logFile   = "distributor_log.json"
	cacheFile = "cache.json"
)

type Distributor struct {
	config              *config.Config
	chanRequests        chan map[string][]mentoring_request.MentoringRequest
	log                 *logrus.Logger
	distributedRequests distributedRequestCache
	slackClient         *slack.Client
}

type distributedRequestCache map[string]mentoring_request.MentoringRequest

func New(cfg *config.Config, chRequests chan map[string][]mentoring_request.MentoringRequest) (*Distributor, error) {
	var d = &Distributor{
		config:              cfg,
		chanRequests:        chRequests,
		log:                 &logrus.Logger{},
		distributedRequests: distributedRequestCache{},
		slackClient:         slack.New(cfg.SlackToken),
	}

	err := logging.SetupLogging(d.log, cfg.LogLevel, logFile)

	err = createCacheFileIfNotExists()
	if err != nil {
		return nil, err
	}

	err = files.JSONToStruct(cacheFile, &d.distributedRequests)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (d *Distributor) Run() {
	for currentMentoringRequests := range d.chanRequests {
		for trackSlug, mentoringRequests := range currentMentoringRequests {
			for _, request := range mentoringRequests {
				if _, alreadySent := d.distributedRequests[request.UUID]; alreadySent {
					continue
				}
				err := d.sendSlackMessage(request, d.config.TrackConfig[trackSlug])
				if err != nil {
					d.log.Error(err)
					continue
				}
				d.log.Info("sent message: ", request.UUID)
				d.distributedRequests[request.UUID] = request
			}
		}
		d.distributedRequests.CleanUp(currentMentoringRequests)

		err := d.distributedRequests.SaveToFile()
		if err != nil {
			d.log.Error(err)
		}
	}

	/*	for requests := range d.chanRequests {
		for _, request := range requests {
			if _, alreadySent := d.distributedRequests[request.UUID]; !alreadySent {
				err := d.sendSlackMessage(request)
				if err != nil {
					d.log.Error(err)
					continue
				}
				d.log.Info("sent message: ", request.UUID)
				d.distributedRequests[request.UUID] = request
			}
		}

		d.distributedRequests.CleanUp(requests)

		err := d.distributedRequests.SaveToFile()
		if err != nil {
			d.log.Error(err)
		}
	}*/
}

func (d Distributor) sendSlackMessage(request mentoring_request.MentoringRequest, trackConfig config.TrackConfig) error {
	attachment := slack.Attachment{
		Pretext: "New mentoring request",
		Text:    fmt.Sprintf("%s: %s", request.UUID, request.URL),
	}

	_, _, err := d.slackClient.PostMessage(
		trackConfig.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionTS(trackConfig.ThreadTS),
	)
	return err
}

func (d distributedRequestCache) CleanUp(currentRequest map[string][]mentoring_request.MentoringRequest) {
outerLoop:
	for _, alreadyDistributedRequest := range d {
		for _, requestsForLanguageTrack := range currentRequest {
			for _, request := range requestsForLanguageTrack {
				if request.UUID == alreadyDistributedRequest.UUID {
					continue outerLoop
				}
			}
		}
		delete(d, alreadyDistributedRequest.UUID)
	}
}

func (d distributedRequestCache) SaveToFile() error {
	file, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheFile, file, 0644)
}

func createCacheFileIfNotExists() error {
	_, err := os.Stat(cacheFile)
	if os.IsNotExist(err) {
		marshal, err := json.Marshal(distributedRequestCache{})
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cacheFile, marshal, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
