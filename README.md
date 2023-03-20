# insights-bot

A bot works with OpenAI GPT models to provide insights for your Telegram info flows.

## Usage

### Run

```shell
OPENAI_API_SECRET=<OpenAI API Secret Key> TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> insights-bot
```

### Run with Docker

```shell
docker run -it --rm -e TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> -e OPENAI_API_SECRET=<OpenAI API Secret Key> insights-bot nekomeowww/insights-bot:latest
```

### Run with docker-compose

Remember to replace your token and cookie in `docker-compose.yml`

```shell
docker-compose up -d
```

## Build on your own

### Build with go

```shell
go build -a -o "release/pero" "github.com/nekomeowww/insights-bot/cmd/pero"
```

### Build with Docker

```shell
docker buildx build --platform linux/arm64,linux/amd64 -t <tag> -f Dockerfile .
```
