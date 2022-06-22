package config

type Config struct {
	LogLevel       string                 `json:"log_level"`
	Interval       int                    `json:"interval"`
	ExercismToken  string                 `json:"exercism_token"`
	SlackToken     string                 `json:"slack_token"`
	TrackConfig    map[string]TrackConfig `json:"track_config"`
	RemindInterval string                 `json:"remind_interval"`
}

type TrackConfig struct {
	ThreadTS  string `json:"thread_ts"`
	ChannelID string `json:"channel_id"`
}
