package configs

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvTelegramBotToken       = "TELEGRAM_BOT_TOKEN" //nolint:gosec
	EnvTelegramBotWebhookURL  = "TELEGRAM_BOT_WEBHOOK_URL"
	EnvTelegramBotWebhookPort = "TELEGRAM_BOT_WEBHOOK_PORT"

	EnvSlackClientID     = "SLACK_CLIENT_ID"
	EnvSlackClientSecret = "SLACK_CLIENT_SECRET"
	EnvSlackWebhookPort  = "SLACK_WEBHOOK_PORT"

	EnvOpenAIAPISecret              = "OPENAI_API_SECRET" //nolint:gosec
	EnvOpenAIAPIHost                = "OPENAI_API_HOST"
	EnvPineconeProjectName          = "PINECONE_PROJECT_NAME"
	EnvPineconeEnvironment          = "PINECONE_ENVIRONMENT"
	EnvPineconeAPIKey               = "PINECONE_API_KEY" //nolint:gosec
	EnvPineconeChatHistoryIndexName = "PINECONE_CHAT_HISTORY_INDEX_NAME"
	EnvDBConnectionString           = "DB_CONNECTION_STR"
)

type SectionPineconeIndexes struct {
	ChatHistoryIndexName string
}

type SectionPinecone struct {
	ProjectName string
	Environment string
	APIKey      string

	Indexes SectionPineconeIndexes
}

type SectionSlack struct {
	Port         string
	ClientID     string
	ClientSecret string
}

type SectionDB struct {
	ConnectionString string
}
type SectionTelegram struct {
	BotToken       string
	BotWebhookURL  string
	BotWebhookPort string
}

type Config struct {
	Telegram        SectionTelegram
	OpenAIAPISecret string
	OpenAIAPIHost   string
	Pinecone        SectionPinecone
	CloverDBPath    string
	DB              SectionDB
	Slack           SectionSlack
}

func NewConfig() func() (*Config, error) {
	return func() (*Config, error) {
		envs, err := godotenv.Read()
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		getEnv := func(varName string) string {
			v, ok := envs[varName]
			if !ok || v == "" {
				return os.Getenv(varName)
			}

			return v
		}

		return &Config{
			Telegram: SectionTelegram{
				BotToken:       getEnv(EnvTelegramBotToken),
				BotWebhookURL:  getEnv(EnvTelegramBotWebhookURL),
				BotWebhookPort: getEnv(EnvTelegramBotWebhookPort),
			},
			Slack: SectionSlack{
				Port:         getEnv(EnvSlackWebhookPort),
				ClientID:     getEnv(EnvSlackClientID),
				ClientSecret: getEnv(EnvSlackClientSecret),
			},
			OpenAIAPISecret: getEnv(EnvOpenAIAPISecret),
			OpenAIAPIHost:   getEnv(EnvOpenAIAPIHost),
			Pinecone: SectionPinecone{
				ProjectName: getEnv(EnvPineconeProjectName),
				Environment: getEnv(EnvPineconeEnvironment),
				APIKey:      getEnv(EnvPineconeAPIKey),
				Indexes: SectionPineconeIndexes{
					ChatHistoryIndexName: getEnv(EnvPineconeChatHistoryIndexName),
				},
			},
			DB: SectionDB{
				ConnectionString: getEnv(EnvDBConnectionString),
			},
		}, nil
	}
}

func NewTestConfig() func() *Config {
	return func() *Config {
		return &Config{
			DB: SectionDB{
				ConnectionString: "postgresql://postgres:123456@localhost:5432/postgres?search_path=public&sslmode=disable",
			},
		}
	}
}
