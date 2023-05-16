<p align="center">
  <image src="./docs/images/icon.png" width="200px" height="200px" />
</p>

<h1 align="center">insights-bot</h1>

<p align="center">
  <img src="https://github.com/nekomeowww/insights-bot/workflows/Testing/badge.svg">
  <img src="https://github.com/nekomeowww/insights-bot/workflows/Building/badge.svg" />
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

---
## Support IMs
- Telegram
- Slack
- Discord

---

## Usage

### Commands

Insights Bot ships with a set of commands, you can use `/help` to get a list of available commands when talking to the bot in Telegram.

#### Summarize webpages

Command: `/smr`

Arguments: URL

Usage:

```txt
/smr https://www.example.com
```

By sending `/smr` command with a URL, the bot will try to summarize the webpage and return the result.

#### Enable chat history recapturing

Command: `/enable_recap`

Arguments: None

```txt
/enable_recap
```

> **Note**
> This command requires user to be an administrator of the chat.

> **Warning**
> **This command will also enable the bot to rapidly send a chat history recap automatically for each 6 hours.**

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

By sending `/enable_recap` command, the bot will start to capture the chat histories and try to summarize them when you send `/recap` command afterwards.

#### Disable chat history recapturing

Command: `/disable_recap`

Arguments: None

```txt
/disable_recap
```

> **Note**
> This command requires user to be an administrator of the chat.

> **Warning**
> **This command will also disable the functionalities of `/recap` command**

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

By sending `/disable_recap` command, the bot will stop capturing the chat histories and no longer respond to `/recap` command.

#### Summarize chat histories or Recap

By sending `/recap` command, the bot will try to summarize the chat histories and return the result you choose later. Such as:

```txt
/recap
```

> **Warning**
> **This command is not available in Slack/Discord integration currently.**

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

### Run with docker-compose

Create your `.env` by making a copy of the contents from `.env.example` file. The .env file should be placed at the root of the project directory next to your `docker-compose.yml` file.

Replace your OpenAI token and other environment variables in `.env`, and then run:

```shell
docker-compose --profile hub up -d
```

If you prefer run docker image from local codes, then run:

```shell
docker-compose --profile local up -d --build
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

## Configurations

### Environment variables

| Name                        | Required    | Default                           | Description                                                                                                                                                                                                                                                                                                                                                                                                             |
|-----------------------------|-------------|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `TELEGRAM_BOT_TOKEN`        | `true`      |                                   | Telegram Bot API token, you can create one and obtain the token through [@BotFather](https://t.me/BotFather)                                                                                                                                                                                                                                                                                                            |
| `TELEGRAM_BOT_WEBHOOK_URL`  | `false`     |                                   | Telegram Bot webhook URL and port, you can use [https://ngrok.com/](https://ngrok.com/) or Cloudflare tunnel to expose your local server to the internet.                                                                                                                                                                                                                                                               |
| `TELEGRAM_BOT_WEBHOOK_PORT` | `false`     | `7071`                            | Telegram Bot Webhook server port, default is 7071                                                                                                                                                                                                                                                                                                                                                                       |
| `OPENAI_API_SECRET`         | `true`      |                                   | OpenAI API Secret Key that looks like `sk-************************************************`, you can obtain one by signing in to OpenAI platform and create one at [http://platform.openai.com/account/api-keys](http://platform.openai.com/account/api-keys).                                                                                                                                                          |
| `OPENAI_API_HOST`           | `false`     | `https://api.openai.com`          | OpenAI API Host, you can specify one if you have a relay or reversed proxy configured. Such as `https://openai.example.workers.dev`                                                                                                                                                                                                                                                                                     |
| `DB_CONNECTION_STR`         | `true`      |                                   | PostgreSQL database URL. Such as `postgres://postgres:postgres@localhost:5432/postgres`. You could also suffix with `?search_path=<schema name>` if you want to specify a schema                                                                                                                                                                                                                                        |
| `SLACK_CLIENT_ID`           | `false`     |                                   | Slack app client id, you can create a slack app and get it, see: [tutorial](https://api.slack.com/tutorials/slack-apps-and-postman)                                                                                                                                                                                                                                                                                     |
| `SLACK_CLIENT_SECRET`       | `false`     |                                   | Slack app client secret, you can create a slack app and get it, see: [tutorial](https://api.slack.com/tutorials/slack-apps-and-postman)                                                                                                                                                                                                                                                                                 |
| `SLACK_WEBHOOK_PORT`        | `false`     | `7070`                            | Port for Slack Bot/App Webhook server, default is 7070                                                                                                                                                                                                                                                                                                                                                                  |
| `DISCORD_BOT_TOKEN`         | `false`     |                                   | Discord bot token, you can create a discord app and get it, see: [Get started document](https://discord.com/developers/docs/getting-started)                                                                                                                                                                                                                                                                            |
| `DISCORD_BOT_PUBLIC_KEY`    | `false`     |                                   | Discord bot public key, you can create a discord app and get it, see: [Get started document](https://discord.com/developers/docs/getting-started), required if `DISCORD_BOT_TOKEN` provided.                                                                                                                                                                                                                            |
| `DISCORD_BOT_WEBHOOK_PORT`  | `false`     | `6060`                            | Port for Discord Bot Webhook server, default is 6060                                                                                                                                                                                                                                                                                                                                                                    |
| ~~`CLOVER_DB_PATH`~~        | ~~`false`~~ | ~~`insights_bot_clover_data.db`~~ | **Deprecated**. ~~Path to Clover database file, you can specify one if you want to specify a path to store data when executed and ran with binary. The default path is `/var/lib/insights-bot/insights_bot_clover_data.db` in Docker volume, you can override the defaults `-e CLOVER_DB_PATH=<path>` when executing `docker run` command or modify and prepend a new `CLOVER_DB_PATH` the `docker-compose.yml` file.~~ |

## Acknowledgements

- Project logo was generated by [Midjourney](https://www.midjourney.com/app/jobs/ff3e9b42-181b-4181-a9ae-6777f957835d/)
- [OpenAI](https://openai.com/) for providing the GPT series models
