package configs

import (
	"os"

	"github.com/samber/lo"
)

const (
	EnvTelegramBotToken             = "TELEGRAM_BOT_TOKEN"
	EnvOpenAIAPISecret              = "OPENAI_API_SECRET"
	EnvPineconeProjectName          = "PINECONE_PROJECT_NAME"
	EnvPineconeEnvironment          = "PINECONE_ENVIRONMENT"
	EnvPineconeAPIKey               = "PINECONE_API_KEY"
	EnvPineconeChatHistoryIndexName = "PINECONE_CHAT_HISTORY_INDEX_NAME"
	EnvCloverDBPath                 = "CLOVER_DB_PATH"
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

type Config struct {
	TelegramBotToken string
	OpenAIAPISecret  string
	Pinecone         SectionPinecone
	CloverDBPath     string
}

func NewConfig() func() *Config {
	return func() *Config {
		return &Config{
			TelegramBotToken: os.Getenv(EnvTelegramBotToken),
			OpenAIAPISecret:  os.Getenv(EnvOpenAIAPISecret),
			Pinecone: SectionPinecone{
				ProjectName: os.Getenv(EnvPineconeProjectName),
				Environment: os.Getenv(EnvPineconeEnvironment),
				APIKey:      os.Getenv(EnvPineconeAPIKey),
				Indexes: SectionPineconeIndexes{
					ChatHistoryIndexName: os.Getenv(EnvPineconeChatHistoryIndexName),
				},
			},
			CloverDBPath: lo.Ternary(os.Getenv(EnvCloverDBPath) != "", os.Getenv(EnvCloverDBPath), "insights_bot_clover_data.db"),
		}
	}
}
