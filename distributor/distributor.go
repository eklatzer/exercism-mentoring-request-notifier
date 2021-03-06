package distributor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/files"
	"exercism-mentoring-request-notifier/request"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

const (
	logFile           = "distributor_log.json"
	exercismColorCode = "#604FCD"
)

//Distributor is used to distribute the infos to Slack
type Distributor struct {
	config              *config.Config
	chanRequests        chan map[string][]request.MentoringRequest
	log                 *logrus.Logger
	distributedRequests distributedRequestCache
	slackClient         slackClient
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

type slackClient interface {
	SendMessage(channel string, options ...slack.MsgOption) (string, string, string, error)
	GetConversationReplies(params *slack.GetConversationRepliesParameters) (msgs []slack.Message, hasMore bool, nextCursor string, err error)
	DeleteMessage(channel, messageTimestamp string) (string, string, error)
}

//New returns an instance of Distributor
func New(cfg *config.Config, chRequests chan map[string][]request.MentoringRequest, cacheFilePath string, setupLogging func(logger *logrus.Logger, level, path string) error) (*Distributor, error) {
	var d = &Distributor{
		config:              cfg,
		chanRequests:        chRequests,
		log:                 &logrus.Logger{},
		distributedRequests: distributedRequestCache{},
		slackClient:         slack.New(cfg.SlackToken),
		cacheFilePath:       cacheFilePath,
	}

	err := setupLogging(d.log, cfg.LogLevel, logFile)
	if err != nil {
		return nil, err
	}

	d.remindInterval, err = time.ParseDuration(cfg.RemindInterval)
	if err != nil {
		return nil, err
	}

	return d, nil
}

//ReadCacheIfExists reads the cache from d.cacheFilePath if exists
func (d *Distributor) ReadCacheIfExists(stat func(name string) (os.FileInfo, error), readFile func(string) ([]byte, error)) error {
	if _, err := stat(d.cacheFilePath); err == nil {
		return files.New(readFile).JSONToStruct(d.cacheFilePath, &d.distributedRequests)
	}
	return nil
}

//Run runs the Distributor and forwards the mentoring requests from the Collector to Slack
func (d *Distributor) Run() {
	for currentMentoringRequests := range d.chanRequests {
		for trackSlug, mentoringRequests := range currentMentoringRequests {
			d.handleRequests(mentoringRequests, trackSlug)
		}

		errors := d.distributedRequests.cleanUp(currentMentoringRequests, d.slackClient)
		for _, err := range errors {
			d.log.Warnf("failed to delete message:%s", err.Error())
		}

		err := d.distributedRequests.saveToFile(d.cacheFilePath)
		if err != nil {
			d.log.Error(err)
		}
	}
}

//StartupCheck checks if all given Slack infos are valid (channel-infos)
func (d Distributor) StartupCheck() error {
	for trackSlug, trackConfig := range d.config.TrackConfig {
		_, _, _, err := d.slackClient.GetConversationReplies(&slack.GetConversationRepliesParameters{
			ChannelID: trackConfig.ChannelID,
			Timestamp: trackConfig.ThreadTS,
			Limit:     1,
		})
		if err != nil {
			return fmt.Errorf("failed to get conversation replies for channel %s, thread %s, track %s:%w", trackConfig.ChannelID, trackConfig.ThreadTS, trackSlug, err)
		}
		time.Sleep(500 * time.Millisecond)
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

		_, messageTimestamp, _, err := d.slackClient.SendMessage(
			d.config.TrackConfig[trackSlug].ChannelID,
			slack.MsgOptionAttachments(slack.Attachment{Text: fmt.Sprintf("%s: <%s|Get to request>", message, req.URL), Fields: []slack.AttachmentField{{Title: "Student", Value: req.StudentHandle}, {Title: "Exercise", Value: req.ExerciseTitle}}, Color: exercismColorCode}),
			slack.MsgOptionTS(d.config.TrackConfig[trackSlug].ThreadTS),
		)
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

func (d distributedRequestCache) cleanUp(currentRequest map[string][]request.MentoringRequest, slackClient slackClient) []error {
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

func (r requestInfo) deleteMessages(client slackClient) []error {
	var errors []error
	for _, message := range r.Messages {
		_, _, err := client.DeleteMessage(message.ChannelID, message.Timestamp)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (d distributedRequestCache) saveToFile(cacheFilePath string) error {
	file, err := json.Marshal(d)
	if err != nil {
		return err
	}
	return os.WriteFile(cacheFilePath, file, 0644)
}
