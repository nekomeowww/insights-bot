package configs

import "os"

const (
	EnvTelegramBotToken = "TELEGRAM_BOT_TOKEN"
	EnvOpenAIAPISecret  = "OPENAI_API_SECRET"
)

type Config struct {
	TelegramBotToken string
	OpenAIAPISecret  string
}

func NewConfig() func() *Config {
	return func() *Config {
		return &Config{
			TelegramBotToken: os.Getenv(EnvTelegramBotToken),
			OpenAIAPISecret:  os.Getenv(EnvOpenAIAPISecret),
		}
	}
}
