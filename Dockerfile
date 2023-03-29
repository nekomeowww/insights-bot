# syntax=docker/dockerfile:1

# 设定构建步骤所使用的来源镜像为基于 Alpine 发行版的 Go 1.19 版本镜像
FROM golang:1.19-alpine as builder

ARG VERSION

# 设定 Go 使用 模块化依赖 管理方式：GO111MODULE
RUN GO111MODULE=on

# 创建路径 /app
RUN mkdir /app

# 复制当前目录下 insights-bot 到 /app/insights-bot
COPY . /app/insights-bot

# 切换到 /app/insights-bot 目录
WORKDIR /app/insights-bot

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go build -a -o "release/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"

# 设定运行步骤所使用的镜像为基于 Alpine 发行版镜像
FROM alpine as runner

# 创建路径 /app
RUN mkdir /app
# 创建路径 /app/insights-bot/bin
RUN mkdir -p /app/insights-bot/bin
# 创建路径 /var/lib/insights-bot
RUN mkdir -p /var/lib/insights-bot
# 配置 CLOVER_DB_PATH 环境变量
ENV CLOVER_DB_PATH /var/lib/insights-bot/insights_bot_clover_data.db

COPY --from=builder /app/insights-bot/release/insights-bot /app/insights-bot/bin/

# 入点是编译好的 hyphen 应用程序
ENTRYPOINT [ "/app/insights-bot/bin/insights-bot" ]
