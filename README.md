# Exercism Mentoring Request Notifier

![test](https://github.com/eklatzer/exercism-mentoring-request-notifier/actions/workflows/test.yml/badge.svg)
![lint](https://github.com/eklatzer/exercism-mentoring-request-notifier/actions/workflows/lint.yml/badge.svg)
![build](https://github.com/eklatzer/exercism-mentoring-request-notifier/actions/workflows/build.yml/badge.svg)
[![codecov](https://codecov.io/gh/eklatzer/exercism-mentoring-request-notifier/branch/master/graph/badge.svg?token=S8X3BI4QCN)](https://codecov.io/gh/eklatzer/exercism-mentoring-request-notifier)

Sends messages to Slack threads when new mentoring requests are created and also creates a reminder if the mentoring request is not accepted after a while.

## Build

`build.sh` can be used to build the current project for multiple platforms. The binaries can then be found at `./bin`. Currently, the tool is build for the following platforms:
* `darwin/386`
* `darwin/amd64`
* `linux/386`
* `linux/amd64`
* `linux/arm`
* `linux/arm64`
* `windows/386`
* `windows/amd64`
* `windows/arm`

### Usage

```console
$ build.sh <package-name>
```

## Usage

````console
exercism-mentoring-request-notifier [flags]

    -cache string
            Defines the path to the cache.json (default "cache.json")
    -config string
            Defines the path to the config (default "config.json")
    -v      bool
            Defines if the version of the current binary should be printed and then exited (default false)
````

## Docker

Use the pre-built image:
```console
$ docker run -d -v /path/to/config/:/go/src/exercism-mentoring-request-notifier/cfg ghcr.io/eklatzer/exercism-mentoring-request-notifier
```

Build the image:

```console
$ docker build -t <tag> .
```

Run container:

````console
$ docker run -d -v /path/to/config/:/go/src/exercism-mentoring-request-notifier/cfg <tag>
````

## Config

````
{
  "log_level": "info",                          //details about log-levels: https://github.com/sirupsen/logrus#level-logging
  "interval": 5,                                //time in seconds after repulling data from the exercism api
  "exercism_token": "<YOUR_EXERCISM_TOKEN>",    //exercism-token: https://exercism.org/settings/api_cli
  "slack_token": "<YOUR_SLACK_TOKEN>",          //https://api.slack.com/tutorials/tracks/getting-a-token
  "remind_interval": "6h",                      //duration after a reminder is created, if the mentoring-requests still exists (possible values: https://pkg.go.dev/maze.io/x/duration#ParseDuration)
  "track_config": {
    "go": {                                     //check out https://exercism.org/api/v2/tracks for further slugs
      "thread_ts": "<SLACK_THREAD_TS>",         //copy the link to the message that should be the start of the thread and extract the ID from the end of the URL (e.g.: /archives/C03L1QZHGRL/p1655583361799859-->1655583361.799859)
      "channel_id": "<SLACK_CHANNEL_ID>"        //channel id of slack-channel: https://help.socialintents.com/article/148-how-to-find-your-slack-team-id-and-slack-channel-id
    },
    "cpp": {
      "thread_ts": "<SLACK_THREAD_TS>",
      "channel_id": "<SLACK_CHANNEL_ID>"
    }
  }
}
````

## Setup

For easier setup use: [setup-mentoring-request-notifier](https://github.com/eklatzer/setup-mentoring-request-notifier)

## Slack requirements

To have the full functionality a Slack token has to be provided. Therefore, an app has to be created: [Create a bot for your workspace](https://slack.com/help/articles/115005265703-Create-a-bot-for-your-workspace) </br>
The needed bot token scopes are:
* `channels:history`: View messages and other content
* `chat:write`: Send messages