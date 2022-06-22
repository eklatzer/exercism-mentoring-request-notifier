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
	logFile           = "distributor_log.json"
	exercismColorCode = "#604FCD"
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

type distributedRequestCache map[string]requestInfo

type requestInfo struct {
	Request  request.MentoringRequest `json:"request"`
	LastSent time.Time                `json:"last_sent"`
	Messages []messageInfo            `json:"messages"`
}

type messageInfo struct {
	ChannelID string `json:"channel_id"`
	Timestamp string `json:"timestamp"`
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

	err = files.New(os.ReadFile).JSONToStruct(d.cacheFilePath, &d.distributedRequests)
	if err != nil {
		return nil, err
	}
	return d, err
}

func (d *Distributor) Run() {
	for currentMentoringRequests := range d.chanRequests {
		for trackSlug, mentoringRequests := range currentMentoringRequests {
			d.handleRequests(mentoringRequests, trackSlug)
		}

		errors := d.distributedRequests.CleanUp(currentMentoringRequests, d.slackClient)
		for _, err := range errors {
			d.log.Warnf("failed to delete message:%s", err.Error())
		}

		err := d.distributedRequests.SaveToFile(d.cacheFilePath)
		if err != nil {
			d.log.Error(err)
		}
	}
}

func (d Distributor) StartupCheck() error {
	for trackSlug, trackConfig := range d.config.TrackConfig {
		_, err := d.sendSlackMessage(trackConfig, slack.Attachment{
			Text:  fmt.Sprintf("Start of mentoring request notifer for `%s`", trackSlug),
			Color: exercismColorCode,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d Distributor) handleRequests(mentoringRequests []request.MentoringRequest, trackSlug string) {
	for _, req := range mentoringRequests {
		info, alreadySent := d.distributedRequests[req.UUID]
		var message = "*New mentoring request*"
		if alreadySent {
			message = "*Reminder*"
		}
		if alreadySent && time.Since(info.LastSent) < d.remindInterval {
			continue
		}
		messageTimestamp, err := d.sendSlackMessage(d.config.TrackConfig[trackSlug], slack.Attachment{Text: fmt.Sprintf("%s: <%s|Get to request>", message, req.URL), Fields: []slack.AttachmentField{{Title: "Student", Value: req.StudentHandle}, {Title: "Exercise", Value: req.ExerciseTitle}}, Color: exercismColorCode})
		if err != nil {
			d.log.Error(err)
			continue
		}

		d.distributedRequests[req.UUID] = requestInfo{
			Request:  req,
			LastSent: time.Now(),
			Messages: append(d.distributedRequests[req.UUID].Messages, messageInfo{ChannelID: d.config.TrackConfig[trackSlug].ChannelID, Timestamp: messageTimestamp}),
		}
	}
}

func (d Distributor) sendSlackMessage(trackConfig config.TrackConfig, attachment slack.Attachment) (string, error) {
	_, messageTimestamp, _, err := d.slackClient.SendMessage(
		trackConfig.ChannelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionTS(trackConfig.ThreadTS),
	)
	return messageTimestamp, err
}

func (d distributedRequestCache) CleanUp(currentRequest map[string][]request.MentoringRequest, slackClient *slack.Client) []error {
	var errors []error
outerLoop:
	for _, alreadyDistributedRequest := range d {
		for _, requestsForLanguageTrack := range currentRequest {
			for _, req := range requestsForLanguageTrack {
				if req.UUID == alreadyDistributedRequest.Request.UUID {
					continue outerLoop
				}
			}
		}
		err := alreadyDistributedRequest.deleteMessages(slackClient)
		errors = append(errors, err...)
		delete(d, alreadyDistributedRequest.Request.UUID)
	}
	return errors
}

func (r requestInfo) deleteMessages(client *slack.Client) []error {
	var errors []error
	for _, message := range r.Messages {
		_, _, err := client.DeleteMessage(message.ChannelID, message.Timestamp)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
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
