package distributor

import (
	"encoding/json"
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/files"
	"exercism-mentoring-request-notifier/logging"
	"exercism-mentoring-request-notifier/request"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"io/ioutil"
	"os"
	"time"
)

const (
	logFile = "distributor_log.json"
)

type Distributor struct {
	config              *config.Config
	chanRequests        chan map[string][]request.MentoringRequest
	log                 *logrus.Logger
	distributedRequests distributedRequestCache
	slackClient         *slack.Client
	cacheFilePath       string
	remindInterval      time.Duration
}

type distributedRequestCache map[string]messageInfo

type messageInfo struct {
	Request  request.MentoringRequest `json:"request"`
	LastSent time.Time                `json:"last_sent"`
}

func New(cfg *config.Config, chRequests chan map[string][]request.MentoringRequest, cacheFilePath string) (*Distributor, error) {
	var d = &Distributor{
		config:              cfg,
		chanRequests:        chRequests,
		log:                 &logrus.Logger{},
		distributedRequests: distributedRequestCache{},
		slackClient:         slack.New(cfg.SlackToken),
		cacheFilePath:       cacheFilePath,
	}

	err := logging.SetupLogging(d.log, cfg.LogLevel, logFile)
	if err != nil {
		return nil, err
	}

	d.remindInterval, err = time.ParseDuration(cfg.RemindInterval)
	if err != nil {
		return nil, err
	}

	err = createCacheFileIfNotExists(cacheFilePath)
	if err != nil {
		return nil, err
	}

	err = files.JSONToStruct(d.cacheFilePath, &d.distributedRequests)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (d *Distributor) Run() {
	for currentMentoringRequests := range d.chanRequests {
		for trackSlug, mentoringRequests := range currentMentoringRequests {
			for _, req := range mentoringRequests {
				info, alreadySent := d.distributedRequests[req.UUID]
				var message = "New mentoring request"
				if alreadySent {
					message = "Reminder"
				}
				if alreadySent && time.Now().Sub(info.LastSent) < d.remindInterval {
					continue
				}
				err := d.sendSlackMessage(req, d.config.TrackConfig[trackSlug], message)
				if err != nil {
					d.log.Error(err)
					continue
				}
				d.log.Info("sent message: ", req.UUID)
				d.distributedRequests[req.UUID] = messageInfo{
					Request:  req,
					LastSent: time.Now(),
				}
			}
		}
		d.distributedRequests.CleanUp(currentMentoringRequests)

		err := d.distributedRequests.SaveToFile(d.cacheFilePath)
		if err != nil {
			d.log.Error(err)
		}
	}
}

func (d Distributor) sendSlackMessage(request request.MentoringRequest, trackConfig config.TrackConfig, message string) error {
	attachment := slack.Attachment{
		Pretext: message,
		Text:    fmt.Sprintf("%s: %s", request.UUID, request.URL),
	}

	_, _, err := d.slackClient.PostMessage(
		trackConfig.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionTS(trackConfig.ThreadTS),
	)
	return err
}

func (d distributedRequestCache) CleanUp(currentRequest map[string][]request.MentoringRequest) {
outerLoop:
	for _, alreadyDistributedRequest := range d {
		for _, requestsForLanguageTrack := range currentRequest {
			for _, req := range requestsForLanguageTrack {
				if req.UUID == alreadyDistributedRequest.Request.UUID {
					continue outerLoop
				}
			}
		}
		delete(d, alreadyDistributedRequest.Request.UUID)
	}
}

func (d distributedRequestCache) SaveToFile(cacheFilePath string) error {
	file, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cacheFilePath, file, 0644)
}

func createCacheFileIfNotExists(cacheFilePath string) error {
	_, err := os.Stat(cacheFilePath)
	if os.IsNotExist(err) {
		marshal, err := json.Marshal(distributedRequestCache{})
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(cacheFilePath, marshal, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
