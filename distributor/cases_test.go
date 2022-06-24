package distributor

import (
	"errors"
	"exercism-mentoring-request-notifier/config"
)

var testCasesNewDistributor = []struct {
	description       string
	expectError       bool
	setupLoggingError error
	config            *config.Config
}{
	{
		description:       "valid setup",
		expectError:       false,
		setupLoggingError: nil,
		config: &config.Config{
			LogLevel:       "info",
			ExercismToken:  "7b35acbf-fb19-4830-8e8f-9e827ce0a86c",
			SlackToken:     "537515cf-90f4-4d92-a56a-5077074938a2",
			TrackConfig:    map[string]config.TrackConfig{},
			RemindInterval: "6h",
		},
	},
	{
		description:       "setup logging fails",
		expectError:       true,
		setupLoggingError: errors.New("test error: setup logging failed"),
		config:            &config.Config{},
	},
	{
		description: "invalid remind interval",
		expectError: true,
		config:      &config.Config{RemindInterval: "invalid value"},
	},
}

var testCasesDistributorReadCacheIfExists = []struct {
	description string
	expectError bool
	statError   error
	readError   error
	content     []byte
}{
	{
		description: "valid cache file",
		expectError: false,
		statError:   nil,
		readError:   nil,
		content:     []byte("{ \"26fd1985-ec74-4e95-8b62-7c2401c8b3bb\": { \"request\": { \"uuid\": \"26fd1985-ec74-4e95-8b62-7c2401c8b3bb\", \"track_title\": \"C++\", \"exercise_icon_url\": \"https://dg8krxphbh767.cloudfront.net/exercises/rna-transcription.svg\", \"exercise_title\": \"Rna Transcription\", \"student_handle\": \"student\", \"student_avatar_url\": \"\", \"updated_at\": \"2022-06-23T20:43:45Z\", \"have_mentored_previously\": true, \"is_favorited\": false, \"status\": null, \"tooltip_url\": \"/api/v2/mentoring/students/onlined?track_slug=cpp\", \"url\": \"https://exercism.org/mentoring/requests/26fd1985-ec74-4e95-8b62-7c2401c8b3bb\" }, \"last_sent\": \"2022-06-24T09:03:32.6500956+02:00\", \"messages\": [ { \"channel_id\": \"722475e2-f38e-11ec-b939-0242ac120002\", \"timestamp\": \"1656054212.105699\" } ] } }"),
	},
	{
		description: "cache file found but read error",
		expectError: true,
		statError:   nil,
		readError:   errors.New("test error: failed to read"),
	},
	{
		description: "no cache file found",
		expectError: false,
		statError:   errors.New("test error: no cache file found"),
	},
}

var testCasesStartupCheck = []struct {
	description                 string
	expectError                 bool
	trackConfig                 map[string]config.TrackConfig
	getConversationRepliesError []error
}{
	{
		description: "no error",
		expectError: false,
		trackConfig: map[string]config.TrackConfig{
			"go":     {ThreadTS: "", ChannelID: ""},
			"cpp":    {ThreadTS: "", ChannelID: ""},
			"csharp": {ThreadTS: "", ChannelID: ""},
			"java":   {ThreadTS: "", ChannelID: ""},
			"php":    {ThreadTS: "", ChannelID: ""},
		},
		getConversationRepliesError: []error{
			nil, nil, nil, nil, nil,
		},
	}, {
		description: "error at third request",
		expectError: true,
		trackConfig: map[string]config.TrackConfig{
			"go":     {ThreadTS: "", ChannelID: ""},
			"cpp":    {ThreadTS: "", ChannelID: ""},
			"csharp": {ThreadTS: "", ChannelID: ""},
			"java":   {ThreadTS: "", ChannelID: ""},
			"php":    {ThreadTS: "", ChannelID: ""},
		},
		getConversationRepliesError: []error{
			nil, nil, errors.New("test error"),
		},
	},
}
