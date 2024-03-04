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

一个能通过与 OpenAI GPT 模型交互来为你的信息流提供洞察的机器人。

[English Documentation](./README.md)

---

## 支持的聊天平台

- Telegram
- Slack
- Discord

---

## 使用方法

### 命令

Insights Bot 附带了一系列的命令，你可以在 Telegram 中与机器人对话时使用 `/help` 来获取可用的命令列表。
如果你遇到了问题，你也可以使用 `/cancel` 来取消任何正在进行的操作。

#### 总结网页

命令：`/smr`

参数: URL（链接），或者回复一条只包含 URL 的消息

用法：

```txt
/smr https://www.example.com
```

```txt
A: https://www.example.com

将 /smr 回复给 A 的消息
```

通过发送 `/smr` 命令并附带一个 URL 或者回复一条只包含 URL 的消息，机器人会尝试总结网页并返回结果。

#### 配置聊天记录回顾

> **Warning**
> **该命令目前不能在 Slack/Discord 平台上使用**

命令：`/configure_recap`

参数：None

```txt
/configure_recap
```

通过在群组中发送 `/configure_recap` 命令，机器人会发送一条消息并包含一些选项，点击按钮来选择你想要配置的项目。

#### 总结聊天记录

> **Warning**
> **该命令目前不能在 Slack/Discord 平台上使用**

命令：`/recap`

参数：None

```txt
/recap
```

通过发送 `/recap` 命令，机器人会尝试总结聊天记录并返回你选择的总结时长范围的结果。

#### 订阅群组聊天记录回顾

> **Warning**
> **该命令目前不能在 Slack/Discord 平台上使用**

命令：`/subscribe_recap`

参数：None

```txt
/subscribe_recap
```

通过发送 `/subscribe_recap` 命令，机器人会开始捕获你订阅的群组的消息，并在可用时通过私聊向你发送一条总结消息的副本。

#### 取消订阅群组聊天记录回顾

> **Warning**
> **该命令目前不能在 Slack/Discord 平台上使用**

命令：`/unsubscribe_recap`

参数：None

```txt
/unsubscribe_recap
```

通过发送 `/unsubscribe_recap` 命令，机器人将不再向你发送订阅的群组的总结消息副本。

#### 总结私聊中转发的消息

> **Warning**
> **该命令目前不能在 Slack/Discord 平台上使用**

命令：`/recap_forwarded_start`, `/recap_forwarded`

参数：None

```txt
/recap_forwarded_start
```

```txt
<转发的消息>
```

```txt
/recap_forwarded
```

通过发送 `/recap_forwarded_start` 命令，机器人会开始捕获你在私聊中转发的消息，并在你发送 `/recap_forwarded` 命令后尝试总结它们。

## 部署

### 使用二进制文件运行

你将会需要克隆这个仓库并且自己构建二进制文件。

```shell
git clone https://github.com/nekomeowww/insights-bot
```

```shell
go build -a -o "build/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"
```

然后将 `.env.example` 文件复制到 `build` 目录并且重命名为 `.env`，然后填写环境变量。

```shell
cd build
cp ../.env.example .env
vim .env
```

```shell
# 配置可执行权限
$ chmod +x ./insights-bot
# 运行
$ ./insights-bot
```

### 使用 Docker 运行

```shell
docker run -it --rm -e TELEGRAM_BOT_TOKEN=<Telegram Bot API 令牌> -e OPENAI_API_SECRET=<OpenAI API 密钥y> -e DB_CONNECTION_STR="<PostgresSQL 连接 URL>" insights-bot ghcr.io/nekomeowww/insights-bot:latest
```

### 使用 Docker Compose 运行

克隆这个项目：

```shell
git clone github.com/nekomeowww/insights-bot
```

或者只复制或下载必要的`.env.example`和`docker-compose.yml`文件（但只能使用预构建的 Docker 镜像来运行 insights-bot）：

```shell
curl -O https://raw.githubusercontent.com/nekomeowww/insights-bot/main/.env.example
curl -O https://raw.githubusercontent.com/nekomeowww/insights-bot/main/docker-compose.yml
```

通过复制 `.env.example` 文件中的内容来创建 `.env` 文件。`.env` 文件应该放在项目根目录下，与 `docker-compose.yml` 文件同级。

```shell
cp .env.example .env
```

通过替换 `.env` 文件中的 OpenAI 令牌和其他环境变量，然后运行：

```shell
docker compose --profile hub up -d
```

如果你想从本地代码编译、构建并运行 Docker 镜像（也就是手动构建，你需要这个项目的全部源代码，可以选择先克隆下来），那么运行：

```shell
docker compose --profile local up -d --build
```

### 亲自构建

#### 使用 Go 构建

```shell
go build -a -o "release/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"
```

#### 使用 Docker 构建

```shell
docker buildx build --platform linux/arm64,linux/amd64 -t <tag> -f Dockerfile .
```

## 所使用的端口

| 端口 | 解释 |
|------|-------------|
| 6060 | pprof Debug 服务器 |
| 7069 | 健康检查服务 |
| 7070 | Slack App/Bot Webhook 服务 |
| 7071 | Telegram Bot Webhook 服务 |
| 7072 | Discord Bot Webhook 服务 |

## 配置

### 环境变量

| 变量名                                           | 必传      | 默认值                                                                                      | 说明                                                                                                                                                                                                                                                                    |
| --------------------------------------------- | ------- | ---------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `TIMEZONE_SHIFT_SECONDS`                      | `false` | `0`                                                                                      | 用于自动生成群组聊天回顾时所参考的时区偏移（以秒为单位），默认值为 0。                                                                                                                                                                                                                                  |
| `TELEGRAM_BOT_TOKEN`                          | `true`  |                                                                                          | Telegram Bot API 令牌，你可以通过 [@BotFather](https://t.me/BotFather) 创建一个。                                                                                                                                                                                                  |
| `TELEGRAM_BOT_WEBHOOK_URL`                    | `false` |                                                                                          | 用于由 Telegram 服务器请求并推送消息更新的 Telegram Bot Webhook URL 以及端口（如果有的话），你可以使用 [https://ngrok.com/](https://ngrok.com/) 或者 Cloudflare tunnel 来讲本地服务暴露到公共互联网。                                                                                                                   |
| `TELEGRAM_BOT_WEBHOOK_PORT`                   | `false` | `7071`                                                                                   | Telegram Bot Webhook 服务监听端口，默认为 7071。                                                                                                                                                                                                                                 |
| `OPENAI_API_SECRET`                           | `true`  |                                                                                          | OpenAI API 密钥，通常类似于 `sk-************************************************` 的结构，你可以登录到 Open AI 并在 [http://platform.openai.com/account/api-keys](http://platform.openai.com/account/api-keys) 上创建一个。                                                                     |
| `OPENAI_API_HOST`                             | `false` | `https://api.openai.com`                                                                 | OpenAI API 的域名，如果配置了中继或反向代理，则可以指定一个。比如 `https://openai.example.workers.dev`                                                                                                                                                                                           |
| `OPENAI_API_MODEL_NAME`                       | `false` | `gpt-3.5-turbo`                                                                          | OpenAI API 模型名称，默认为 `gpt-3.5-turbo`，如果你使用其他模型，比如  `gpt-4` 则可以制指定一个。                                                                                                                                                                                                   |
| `OPENAI_API_TOKEN_LIMIT`                      | `false` | `4096`                                                                                   | OpenAI API Token 限制，用于在调用 Chat Completion API 之前计算文本的分割和截断，一般设置为模型的最大令牌限制，然后交由 insights-bot 决定如何处理，默认为 `4096`。                                                                                                                                                        |
| `OPENAI_API_CHAT_HISTORIES_RECAP_TOKEN_LIMIT` | `false` | `2000`                                                                                   | OpenAI 聊天历史记录回顾令牌限制，生成的和响应的聊天历史记录回顾消息的令牌长度，默认值为 2000，这将会给实际的聊天上下文留下 `OPENAI_API_TOKEN_LIMIT` - 2000 个令牌                                                                                                                                                               |
| `DB_CONNECTION_STR`                           | `true`  | `postgresql://postgres:123456@db_local:5432/postgres?search_path=public&sslmode=disable` | PostgreSQL 数据库连接 URL。结构类似于 `postgres://postgres:postgres@localhost:5432/postgres`。如果你需要指定 schema，则可以通过在后缀加上 `?search_path=<schema name>` 来实现。                                                                                                                         |
| `SLACK_CLIENT_ID`                             | `false` |                                                                                          | Slack app client id，你可以参考[教程](https://api.slack.com/tutorials/slack-apps-and-postman)来创建一个。                                                                                                                                                                           |
| `SLACK_CLIENT_SECRET`                         | `false` |                                                                                          | Slack app client secret，你可以参考[教程](https://api.slack.com/tutorials/slack-apps-and-postman)来创建一个。                                                                                                                                                                       |
| `SLACK_WEBHOOK_PORT`                          | `false` | `7070`                                                                                   | Slack Bot/App Webhook 服务监听端口，默认为 7070。                                                                                                                                                                                                                                |
| `DISCORD_BOT_TOKEN`                           | `false` |                                                                                          | Discord bot token，你可以通过参考[快速上手指南](https://discord.com/developers/docs/getting-started)来创建一个。                                                                                                                                                                          |
| `DISCORD_BOT_PUBLIC_KEY`                      | `false` |                                                                                          | Discord bot public key，你可以通过参考[快速上手指南](https://discord.com/developers/docs/getting-started)来创建一个。当 `DISCORD_BOT_TOKEN` 传入时必传。                                                                                                                                         |
| `DISCORD_BOT_WEBHOOK_PORT`                    | `false` | `7072`                                                                                   | Discord Bot Webhook 服务监听端口，默认为 7072。                                                                                                                                                                                                                                  |
| `REDIS_HOST`                                  | `true`  | `localhost`                                                                              | Redis 主机，默认为 `localhost`                                                                                                                                                                                                                                              |
| `REDIS_PORT`                                  | `true`  | `6379`                                                                                   | Redis 端口，默认为  `6379`                                                                                                                                                                                                                                                  |
| `REDIS_TLS_ENABLED`                           | `false` | `false`                                                                                  | 是否开启 Redis TLS，默认为  `false`                                                                                                                                                                                                                                           |
| `REDIS_USERNAME`                              | `false` |                                                                                          | Redis 用户名。                                                                                                                                                                                                                                                            |
| `REDIS_PASSWORD`                              | `false` |                                                                                          | Redis 密码。                                                                                                                                                                                                                                                             |
| `REDIS_DB`                                    | `false` | `0`                                                                                      | Redis 数据库，默认为 `0`                                                                                                                                                                                                                                                     |
| `REDIS_CLIENT_CACHE_ENABLED`                  | `false` | `false`                                                                                  | 是否开启 Redis 客户端缓存，默认为 `false`, 你可以在这里了解更多相关的知识：[https://redis.io/docs/manual/client-side-caching/](https://redis.io/docs/manual/client-side-caching/) 以及 [https://github.com/redis/rueidis#client-side-caching](https://github.com/redis/rueidis#client-side-caching)。 |
| `LOG_FILE_PATH`                               | `false` | `<insights-bot_executable>/logs/insights-bot.log`                                        | 日志文件路径，如果你想指定二进制执行和运行时存储日志的路径，可以指定一个。默认路径是 Docker 卷中的 `/var/log/insights-bot/insights-bot.log`，你可以在执行 `docker run` 命令时覆盖默认路径 `-e LOG_FILE_PATH=<path>` 或修改并在 `docker-compose.yml` 文件中预置新的 `LOG_FILE_PATH` 。                                                           |
| `LOG_LEVEL`                                   | `false` | `info`                                                                                   | 日志等级，可选值为 `debug`，`info`，`warn`， `error`。                                                                                                                                                                                                                             |
| `LOCALES_DIR`                                 | `false` | `locales`                                                                                | 本地化目录，默认值为 `locales`，推荐配置为绝对路径。                                                                                                                                                                                                                              |

## 鸣谢

- 项目图标由 [Midjourney](https://www.midjourney.com/app/jobs/ff3e9b42-181b-4181-a9ae-6777f957835d/) 生成。
- [OpenAI](https://openai.com/) 提供了强大的 GPT 模型。
