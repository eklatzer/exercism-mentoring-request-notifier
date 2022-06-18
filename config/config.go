package config

import (
	"exercism-mentoring-request-notifier/files"
)

type Config struct {
	TrackSlug     string `json:"track_slug"`
	LogLevel      string `json:"log_level"`
	Interval      int    `json:"interval"`
	ExercismToken string `json:"exercism_token"`
	ChannelID     string `json:"channel_id"`
	SlackToken    string `json:"slack_token"`
	ThreadTS      string `json:"thread_ts"`
}

func ReadConfig(path string) (*Config, error) {
	c := &Config{}
	err := files.JSONToStruct(path, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
