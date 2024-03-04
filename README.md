<p align="center">
  <image src="./docs/images/icon.png" width="200px" height="200px" />
</p>

<h1 align="center">insights-bot</h1>

<p align="center">
  <img src="https://github.com/nekomeowww/insights-bot/actions/workflows/ci.yml/badge.svg">
  <img src="https://github.com/nekomeowww/insights-bot/actions/workflows/build.yml/badge.svg" />
  <img src="https://github.com/nekomeowww/insights-bot/actions/workflows/release-build.yml/badge.svg" />
  <a href="https://goreportcard.com/badge/github.com/nekomeowww/insights-bot"><img src="https://goreportcard.com/badge/github.com/nekomeowww/insights-bot" /></a>
  <a href="https://hub.docker.com/r/nekomeowww/insights-bot">
    <img src="https://img.shields.io/docker/pulls/nekomeowww/insights-bot" />
  </a>
  <a href="https://hub.docker.com/r/nekomeowww/insights-bot">
    <img src="https://img.shields.io/docker/v/nekomeowww/insights-bot" />
  </a>
</p>

<p align="center">
  <a href="https://t.me/ayaka_insights_bot_group">
    <img src="https://img.shields.io/badge/Chat%20on-Telegram-%235AA9E6?logo=telegram" />
  </a>
  <a href="https://slack.com/oauth/v2/authorize?client_id=2628877438886.5144808095409&scope=chat:write,commands&user_scope=">
    <img src="https://img.shields.io/badge/Add_to_Slack-4A154B?logo=slack" />
  </a>
</p>

A bot works with OpenAI GPT models to provide insights for your info flows.

[简体中文](./README.zh-CN.md)

---

## Supported IMs

- Telegram
- Slack
- Discord

---

## Usage

### Commands

Insights Bot ships with a set of commands, you can use `/help` to get a list of available commands when talking to the bot in Telegram.
You can also use `/cancel` to cancel any ongoing actions with the bot.

#### Summarize webpages

Command: `/smr`

Arguments: URL, Replied message with only URL

Usage:

```txt
/smr https://www.example.com
```

```txt
/smr [Reply to a message with only URL]
```

By sending `/smr` command with a URL or replying to a message that only contains a URL, the bot will try to summarize the webpage and return the result.

#### Configure chat history recapturing

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

Command: `/configure_recap`

Arguments: None

```txt
/configure_recap
```

By sending `/configure_recap` command, the bot will send you a message with options you can interact with. Click the buttons to choose the settings you want to configure.

#### Summarize chat histories or Recap

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

Command: `/recap`

Arguments: None

```txt
/recap
```

By sending `/recap` command, the bot will try to summarize the chat histories and return the result you choose later.

#### Subscribe to chat histories recap for a group

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

Command: `/subscribe_recap`

Arguments: None

```txt
/subscribe_recap
```

By sending `/subscribe_recap` command, the bot will start to capture the messages from the group you subscribed and then send a copy of the recap message to you through private chat when it is available.

#### Unsubscribe to chat histories recap for a group

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

Command: `/unsubscribe_recap`

Arguments: None

```txt
/unsubscribe_recap
```

By sending `/unsubscribe_recap` command, the bot will no longer send the copy of the recap message for the group you subscribe.

#### Summarize forwarded messages in private chat

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

Commands: `/recap_forwarded_start`, `/recap_forwarded`

Arguments: None

```txt
/recap_forwarded_start
```

```txt
<Forwarded messages>
```

```txt
/recap_forwarded
```

By sending `/recap_forwarded_start` command, the bot will start to capture the forwarded messages you send later in private chat and try to summarize them when you send `/recap_forwarded` command afterwards.

## Deployment

### Run with binary

You will have to clone this repository and then build the binary by yourself.

```shell
git clone https://github.com/nekomeowww/insights-bot
```

```shell
go build -a -o "build/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"
```

Then copy the `.env.example` file to `build` directory and rename it to `.env`, and then fill in the environment variables.

```shell
cd build
cp ../.env.example .env
vim .env
```

```shell
# assign executable permission to the binary
$ chmod +x ./insights-bot
# run the binary
$ ./insights-bot
```

### Run with docker

```shell
docker run -it --rm -e TELEGRAM_BOT_TOKEN=<Telegram Bot API Token> -e OPENAI_API_SECRET=<OpenAI API Secret Key> -e DB_CONNECTION_STR="<PostgresSQL connection URL>" insights-bot ghcr.io/nekomeowww/insights-bot:latest
```

### Run with Docker Compose

Clone this project:

```shell
git clone github.com/nekomeowww/insights-bot
```

Or only copy or download the necessary `.env.example` and `docker-compose.yml` files (but you will only be able to run the bot with pre-bundled docker image):

```shell
curl -O https://raw.githubusercontent.com/nekomeowww/insights-bot/main/.env.example
curl -O https://raw.githubusercontent.com/nekomeowww/insights-bot/main/docker-compose.yml
```

Create your `.env` by making a copy of the contents from `.env.example` file. The .env file should be placed at the root of the project directory next to your `docker-compose.yml` file.

```shell
cp .env.example .env
```

Replace your OpenAI token and other environment variables in `.env`, and then run:

```shell
docker compose --profile hub up -d
```

If you prefer run docker image from local codes (which means build it manually, you will need the entire source code of this project), then run:

```shell
docker compose --profile local up -d --build
```

### Build on your own

#### Build with go

```shell
go build -a -o "release/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"
```

#### Build with Docker

```shell
docker buildx build --platform linux/arm64,linux/amd64 -t <tag> -f Dockerfile .
```

## Ports we use

| Port | Description |
|------|-------------|
| 6060 | pprof Debug server |
| 7069 | Health check server |
| 7070 | Slack App/Bot webhook server |
| 7071 | Telegram Bot webhook server |
| 7072 | Discord Bot webhook server |

## Configurations

### Environment variables

| Name                                          | Required | Default                                                                                  | Description                                                                                                                                                                                                                                                                                                                                                             |
| --------------------------------------------- | -------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `TIMEZONE_SHIFT_SECONDS`                      | `false`  | `0`                                                                                      | Timezone shift in seconds used for auto generating recap messages for groups, default is 0.                                                                                                                                                                                                                                                                             |
| `TELEGRAM_BOT_TOKEN`                          | `true`   |                                                                                          | Telegram Bot API token, you can create one and obtain the token through [@BotFather](https://t.me/BotFather)                                                                                                                                                                                                                                                            |
| `TELEGRAM_BOT_WEBHOOK_URL`                    | `false`  |                                                                                          | Telegram Bot webhook URL and port, you can use [https://ngrok.com/](https://ngrok.com/) or Cloudflare tunnel to expose your local server to the internet.                                                                                                                                                                                                               |
| `TELEGRAM_BOT_WEBHOOK_PORT`                   | `false`  | `7071`                                                                                   | Telegram Bot Webhook server port, default is 7071                                                                                                                                                                                                                                                                                                                       |
| `OPENAI_API_SECRET`                           | `true`   |                                                                                          | OpenAI API Secret Key that looks like `sk-************************************************`, you can obtain one by signing in to OpenAI platform and create one at [http://platform.openai.com/account/api-keys](http://platform.openai.com/account/api-keys).                                                                                                          |
| `OPENAI_API_HOST`                             | `false`  | `https://api.openai.com`                                                                 | OpenAI API Host, you can specify one if you have a relay or reversed proxy configured. Such as `https://openai.example.workers.dev`                                                                                                                                                                                                                                     |
| `OPENAI_API_MODEL_NAME`                       | `false`  | `gpt-3.5-turbo`                                                                          | OpenAI API model name, default is `gpt-3.5-turbo`, you can specify one if you want to use another model. Such as `gpt-4`                                                                                                                                                                                                                                                |
| `OPENAI_API_TOKEN_LIMIT`                      | `false`  | `4096`                                                                                   | OpenAI API token limit used to computed the splits and truncations of texts before calling Chat Completion API generally set to the maximum token limit of a model, and let insights-bot to determine how to process it, default is `4096`                                                                                                                              |
| `OPENAI_API_CHAT_HISTORIES_RECAP_TOKEN_LIMIT` | `false`  | `2000`                                                                                   | OpenAI chat histories recap token limit, token length of generated and response chat histories recap message, default is 2000, this will leave OPENAI_API_TOKEN_LIMIT - 2000 tokens for actual chat context.                                                                                                                                                            |
| `DB_CONNECTION_STR`                           | `true`   | `postgresql://postgres:123456@db_local:5432/postgres?search_path=public&sslmode=disable` | PostgreSQL database URL. Such as `postgres://postgres:postgres@localhost:5432/postgres`. You could also suffix with `?search_path=<schema name>` if you want to specify a schema.                                                                                                                                                                                       |
| `SLACK_CLIENT_ID`                             | `false`  |                                                                                          | Slack app client id, you can create a slack app and get it, see: [tutorial](https://api.slack.com/tutorials/slack-apps-and-postman)                                                                                                                                                                                                                                     |
| `SLACK_CLIENT_SECRET`                         | `false`  |                                                                                          | Slack app client secret, you can create a slack app and get it, see: [tutorial](https://api.slack.com/tutorials/slack-apps-and-postman)                                                                                                                                                                                                                                 |
| `SLACK_WEBHOOK_PORT`                          | `false`  | `7070`                                                                                   | Port for Slack Bot/App Webhook server, default is 7070                                                                                                                                                                                                                                                                                                                  |
| `DISCORD_BOT_TOKEN`                           | `false`  |                                                                                          | Discord bot token, you can create a discord app and get it, see: [Get started document](https://discord.com/developers/docs/getting-started)                                                                                                                                                                                                                            |
| `DISCORD_BOT_PUBLIC_KEY`                      | `false`  |                                                                                          | Discord bot public key, you can create a discord app and get it, see: [Get started document](https://discord.com/developers/docs/getting-started), required if `DISCORD_BOT_TOKEN` provided.                                                                                                                                                                            |
| `DISCORD_BOT_WEBHOOK_PORT`                    | `false`  | `7072`                                                                                   | Port for Discord Bot Webhook server, default is 7702                                                                                                                                                                                                                                                                                                                    |
| `REDIS_HOST`                                  | `true`   | `localhost`                                                                              | Redis host connects to, default is `localhost`                                                                                                                                                                                                                                                                                                                          |
| `REDIS_PORT`                                  | `true`   | `6379`                                                                                   | Redis port, default is `6379`                                                                                                                                                                                                                                                                                                                                           |
| `REDIS_TLS_ENABLED`                           | `false`  | `false`                                                                                  | Redis TLS enabled, default is `false`                                                                                                                                                                                                                                                                                                                                   |
| `REDIS_USERNAME`                              | `false`  |                                                                                          | Redis username.                                                                                                                                                                                                                                                                                                                                                         |
| `REDIS_PASSWORD`                              | `false`  |                                                                                          | Redis password.                                                                                                                                                                                                                                                                                                                                                         |
| `REDIS_DB`                                    | `false`  | `0`                                                                                      | Redis database, default is `0`                                                                                                                                                                                                                                                                                                                                          |
| `REDIS_CLIENT_CACHE_ENABLED`                  | `false`  | `false`                                                                                  | Redis client cache enabled, default is `false`, read more about client cache at [https://redis.io/docs/manual/client-side-caching/](https://redis.io/docs/manual/client-side-caching/) and [https://github.com/redis/rueidis#client-side-caching](https://github.com/redis/rueidis#client-side-caching) for more details.                                               |
| `LOG_FILE_PATH`                               | `false`  | `<insights-bot_executable>/logs/insights-bot.log`                                        | Log file path, you can specify one if you want to specify a path to store logs when executed and ran with binary. The default path is `/var/log/insights-bot/insights-bot.log` in Docker volume, you can override the defaults `-e LOG_FILE_PATH=<path>` when executing `docker run` command or modify and prepend a new `LOG_FILE_PATH` the `docker-compose.yml` file. |
| `LOG_LEVEL`                                   | `false`  | `info`                                                                                   | Log level, available values are `debug`, `info`, `warn`, `error`                                                                                                                                                                                                                                                                                                        |
| `LOCALES_DIR`                                 | `false`  | `locales`                                                                                | Locales directory, default is `locales`, it is recommended to configure as an absolute path.                                                                                                                                                                                                                                                                            |

## Acknowledgements

- Project logo was generated by [Midjourney](https://www.midjourney.com/app/jobs/ff3e9b42-181b-4181-a9ae-6777f957835d/)
- [OpenAI](https://openai.com/) for providing the GPT series models
