<p align="center">
  <img width=300px src="https://github.com/seitau/alfred-crypto-unit-converter/assets/32608705/bb507f96-b429-4331-9d56-25f4ccfb8091" />
</p>

<p align="center">
 Bilingual Cat is your purrfectly fluffy assistant.
</p>

# Bilingual Cat

## Build

```sh
make build
```

## Docker Build

```sh
make docker-build
```

## Run

```sh
SLACK_BOT_TOKEN={slack bot token} \
SLACK_WEBHOOK_SIGNING_SECRET={slack signing secret} \
OPEN_AI_API_TOKEN={open ai api token} \
go run main.go
```

Checkout [slack apps page](https://api.slack.com/apps/A057JFH0E5T) to get secret and bot token, and configure webhook endpoint.

- webhook url configuration: https://api.slack.com/apps/A057JFH0E5T/event-subscriptions
- signing secret: https://api.slack.com/apps/A057JFH0E5T/general
  - signing secret is used to verify webhook payload sent from slack.
- bot token: https://api.slack.com/apps/A057JFH0E5T/oauth
  - bot token begins with `xoxb`
