package collector

import (
	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/request"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	var cfg = &config.Config{LogLevel: "warn", Interval: 5, ExercismToken: "eb32823f-b97e-4480-9737-e2db5cd660db", SlackToken: "22d06021-b19b-4b61-ba98-f41df43392ba", TrackConfig: map[string]config.TrackConfig{}, RemindInterval: "6h"}
	var ch = make(chan map[string][]request.MentoringRequest)
	collector, err := New(cfg, ch, func(logger *logrus.Logger, level, path string) error { return nil })
	assert.Nil(t, err)
	assert.Equal(t, cfg, collector.config)
	assert.Equal(t, ch, collector.chanRequests)
}
