# syntax=docker/dockerfile:1

# --- builder ---
FROM golang:1.21 as builder

RUN GO111MODULE=on

RUN mkdir /app

COPY . /app/insights-bot

WORKDIR /app/insights-bot

RUN go env
RUN go env -w CGO_ENABLED=0
RUN go mod download
RUN go build -a -o "release/insights-bot" "github.com/nekomeowww/insights-bot/cmd/insights-bot"

# --- runner ---
FROM debian as runner

RUN apt update && apt upgrade -y && apt install -y ca-certificates curl && update-ca-certificates

COPY --from=builder /app/insights-bot/release/insights-bot /usr/local/bin/
COPY --from=builder /app/insights-bot/locales /etc/insights-bot/locales

ENV LOG_FILE_PATH /var/log/insights-bot/insights-bot.log
ENV LOCALES_DIR /etc/insights-bot/locales

EXPOSE 7069

CMD [ "/usr/local/bin/insights-bot" ]
