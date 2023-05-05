package configs

import (
	"os"
)

const (
	EnvTelegramBotToken = "TELEGRAM_BOT_TOKEN" //nolint:gosec

	EnvSlackClientID     = "SLACK_CLIENT_ID"
	EnvSlackClientSecret = "SLACK_CLIENT_SECRET"
	EnvSlackBotPort      = "SLACK_BOT_PORT"

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

type Config struct {
	TelegramBotToken string
	OpenAIAPISecret  string
	OpenAIAPIHost    string
	Pinecone         SectionPinecone
	CloverDBPath     string
	DB               SectionDB
	Slack            SectionSlack
}

func NewConfig() func() *Config {
	return func() *Config {
		return &Config{
			TelegramBotToken: os.Getenv(EnvTelegramBotToken),
			Slack: SectionSlack{
				Port:         os.Getenv(EnvSlackBotPort),
				ClientID:     os.Getenv(EnvSlackClientID),
				ClientSecret: os.Getenv(EnvSlackClientSecret),
			},
			OpenAIAPISecret: os.Getenv(EnvOpenAIAPISecret),
			OpenAIAPIHost:   os.Getenv(EnvOpenAIAPIHost),
			Pinecone: SectionPinecone{
				ProjectName: os.Getenv(EnvPineconeProjectName),
				Environment: os.Getenv(EnvPineconeEnvironment),
				APIKey:      os.Getenv(EnvPineconeAPIKey),
				Indexes: SectionPineconeIndexes{
					ChatHistoryIndexName: os.Getenv(EnvPineconeChatHistoryIndexName),
				},
			},
			DB: SectionDB{
				ConnectionString: os.Getenv(EnvDBConnectionString),
			},
		}
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
