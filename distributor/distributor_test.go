package distributor

import (
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/request"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	for _, testCase := range testCasesNewDistributor {
		t.Run(testCase.description, func(t *testing.T) {
			var ch = make(chan map[string][]request.MentoringRequest)
			distributor, err := New(testCase.config, ch, "", func(logger *logrus.Logger, level, path string) error {
				return testCase.setupLoggingError
			})
			assertError(t, err, testCase.expectError)
			if testCase.expectError {
				assert.Nil(t, distributor)
			} else {
				assert.NotNil(t, distributor)
			}
		})
	}
}

func TestDistributor_ReadCacheIfExists(t *testing.T) {
	for _, testCase := range testCasesDistributorReadCacheIfExists {
		t.Run(testCase.description, func(t *testing.T) {
			var d = &Distributor{cacheFilePath: "cache.json"}
			err := d.ReadCacheIfExists(func(name string) (os.FileInfo, error) {
				return nil, testCase.statError
			}, func(s string) ([]byte, error) {
				return testCase.content, testCase.readError
			})
			assertError(t, err, testCase.expectError)
			if len(testCase.content) != 0 {
				assert.NotNil(t, d.distributedRequests)
			}
		})
	}
}

type mockSlackClient struct {
	sendMessageError               []error
	getConversationRepliesError    []error
	deleteMessageError             []error
	numberOfSentMessages           int
	numberOfGetConversationReplies int
	numberOfDeleteMessage          int
}

func (m *mockSlackClient) SendMessage(channel string, options ...slack.MsgOption) (string, string, string, error) {
	m.numberOfSentMessages++
	return "", "", "", m.sendMessageError[m.numberOfSentMessages-1]
}

func (m *mockSlackClient) GetConversationReplies(params *slack.GetConversationRepliesParameters) (msgs []slack.Message, hasMore bool, nextCursor string, err error) {
	m.numberOfGetConversationReplies++
	return nil, false, "", m.getConversationRepliesError[m.numberOfGetConversationReplies-1]
}

func (m *mockSlackClient) DeleteMessage(channel, messageTimestamp string) (string, string, error) {
	m.numberOfDeleteMessage++
	return "", "", m.deleteMessageError[m.numberOfDeleteMessage-1]
}

func TestDistributor_StartupCheck(t *testing.T) {
	for _, testCase := range testCasesStartupCheck {
		t.Run(testCase.description, func(t *testing.T) {
			var d = &Distributor{
				slackClient: &mockSlackClient{
					getConversationRepliesError: testCase.getConversationRepliesError,
				},
				config: &config.Config{TrackConfig: testCase.trackConfig},
			}
			err := d.StartupCheck()
			assertError(t, err, testCase.expectError)
		})
	}
}

func assertError(t *testing.T, err error, expectError bool) {
	if expectError {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}
}
