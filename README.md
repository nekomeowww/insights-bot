# insights-bot

A bot works with OpenAI GPT models to provide insights for your Telegram info flows.

## Usage

### Run

```shell
OPENAI_API_SECRET=<OpenAI API Secret Key> TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> CLOVER_DB_PATH=<path to store DB> insights-bot
```

#### Run with a specific OpenAI API endpoint host

```shell
OPENAI_API_HOST=https://<Some Host> OPENAI_API_SECRET=<OpenAI API Secret Key> TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> CLOVER_DB_PATH=<path to store DB> insights-bot
```

### Run with Docker

```shell
docker run -it --rm -e TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> -e OPENAI_API_SECRET=<OpenAI API Secret Key> -e CLOVER_DB_PATH=<path to store DB> insights-bot nekomeowww/insights-bot:latest
```

#### Run with docker and a specific OpenAI API endpoint host

```shell
docker run -it --rm -e TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> -e OPENAI_API_HOST=https://<Some Host> -e OPENAI_API_SECRET=<OpenAI API Secret Key> -e CLOVER_DB_PATH=<path to store DB> insights-bot nekomeowww/insights-bot:latest
```

### Run with docker-compose

Remember to replace your token and cookie in `docker-compose.yml`

```shell
docker-compose up -d
```

If you prefer run docker image from local codes,

```shell
docker-compose --profile local up -d --build
```

## Build on your own

### Build with go

```shell
go build -a -o "release/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"
```

### Build with Docker

```shell
docker buildx build --platform linux/arm64,linux/amd64 -t <tag> -f Dockerfile .
```
