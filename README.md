# Exercism Mentoring Request Notifier

Sends messages to Slack threads when new mentoring requests are created for a certain track.

## Usage

````
exercism-mentoring-request-notifier [flags]

  -config string
        Defines the path to the config (default "config.json")
````

## Config

````json
{
  "track_slug": "go",                         //check out https://exercism.org/api/v2/tracks for further slugs
  "log_level": "info",                        //details about log-levels: https://github.com/sirupsen/logrus#level-logging
  "interval": 5,                              //time in seconds after repulling data from the exercism api
  "exercism_token": "<YOUR_EXERCISM_TOKEN>",  //exercism-token: https://exercism.org/settings/api_cli
  "channel_id": "<SLACK_CHANNEL_ID>",         //channel id of slack-channel: https://help.socialintents.com/article/148-how-to-find-your-slack-team-id-and-slack-channel-id
  "slack_token": "<YOUR_SLACK_TOKEN>",        //https://api.slack.com/tutorials/tracks/getting-a-token
  "thread_ts": "<SLACK_THREAD_TS>"            //copy the link to the message that should be the start of the thread and extract the ID from the end of the URL (e.g.: /archives/C03L1QZHGRL/p1655583361799859-->1655583361.799859)
}
````