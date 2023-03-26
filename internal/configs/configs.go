package configs

import "os"

const (
	EnvTelegramBotToken    = "TELEGRAM_BOT_TOKEN"
	EnvOpenAIAPISecret     = "OPENAI_API_SECRET"
	EnvPineconeIndexName   = "PINECONE_INDEX_NAME"
	EnvPineconeProjectName = "PINECONE_PROJECT_NAME"
	EnvPineconeEnvironment = "PINECONE_ENVIRONMENT"
	EnvPineconeAPIKey      = "PINECONE_API_KEY"
)

type SectionPinecone struct {
	IndexName   string
	ProjectName string
	Environment string
	APIKey      string
}

type Config struct {
	TelegramBotToken string
	OpenAIAPISecret  string
	Pinecone         SectionPinecone
}

func NewConfig() func() *Config {
	return func() *Config {
		return &Config{
			TelegramBotToken: os.Getenv(EnvTelegramBotToken),
			OpenAIAPISecret:  os.Getenv(EnvOpenAIAPISecret),
			Pinecone: SectionPinecone{
				IndexName:   os.Getenv(EnvPineconeIndexName),
				ProjectName: os.Getenv(EnvPineconeProjectName),
				Environment: os.Getenv(EnvPineconeEnvironment),
				APIKey:      os.Getenv(EnvPineconeAPIKey),
			},
		}
	}
}
