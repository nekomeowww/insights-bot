package configs

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/samber/lo"
	goopenai "github.com/sashabaranov/go-openai"
)

const (
	EnvTimezoneShiftSeconds = "TIMEZONE_SHIFT_SECONDS"

	EnvTelegramBotToken       = "TELEGRAM_BOT_TOKEN" //nolint:gosec
	EnvTelegramBotWebhookURL  = "TELEGRAM_BOT_WEBHOOK_URL"
	EnvTelegramBotWebhookPort = "TELEGRAM_BOT_WEBHOOK_PORT"
	EnvTelegramBotAPIEndpoint = "TELEGRAM_BOT_API_ENDPOINT"

	EnvSlackClientID     = "SLACK_CLIENT_ID"
	EnvSlackClientSecret = "SLACK_CLIENT_SECRET"
	EnvSlackWebhookPort  = "SLACK_WEBHOOK_PORT"

	EnvDiscordBotToken       = "DISCORD_BOT_TOKEN" //nolint:gosec
	EnvDiscordBotPublicKey   = "DISCORD_BOT_PUBLIC_KEY"
	EnvDiscordBotWebhookPort = "DISCORD_BOT_WEBHOOK_PORT"

	EnvOpenAIAPISecret    = "OPENAI_API_SECRET" //nolint:gosec
	EnvOpenAIAPIHost      = "OPENAI_API_HOST"
	EnvOpenAIAPIModelName = "OPENAI_API_MODEL_NAME"

	EnvPineconeProjectName          = "PINECONE_PROJECT_NAME"
	EnvPineconeEnvironment          = "PINECONE_ENVIRONMENT"
	EnvPineconeAPIKey               = "PINECONE_API_KEY" //nolint:gosec
	EnvPineconeChatHistoryIndexName = "PINECONE_CHAT_HISTORY_INDEX_NAME"

	EnvDBConnectionString = "DB_CONNECTION_STR"

	EnvRedisHost               = "REDIS_HOST"
	EnvRedisPort               = "REDIS_PORT"
	EnvRedisTLSEnabled         = "REDIS_TLS_ENABLED"
	EnvRedisUsername           = "REDIS_USERNAME"
	EnvRedisPassword           = "REDIS_PASSWORD"
	EnvRedisDB                 = "REDIS_DB"
	EnvRedisClientCacheEnabled = "REDIS_CLIENT_CACHE_ENABLED"

	EnvLogLevel    = "LOG_LEVEL"
	EnvLogFilePath = "LOG_FILE_PATH"

	EnvHardLimitManualRecapRatePerSeconds      = "HARD_LIMIT_MANUAL_RECAP_RATE_PER_SECONDS"
	EnvHardLimitSummarizeWebpageRatePerSeconds = "HARD_LIMIT_SMR_WEBPAGE_RATE_PER_SECONDS"
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

type SectionDiscord struct {
	Port      string
	Token     string
	PublicKey string
}

type SectionDB struct {
	ConnectionString string
}
type SectionTelegram struct {
	BotToken       string
	BotWebhookURL  string
	BotWebhookPort string
	BotAPIEndpoint string
}

type SectionRedis struct {
	Host               string
	Port               string
	TLSEnabled         bool
	Username           string
	Password           string
	DB                 int64
	ClientCacheEnabled bool
}

type SectionHardLimit struct {
	ManualRecapRatePerSeconds      int64
	SummarizeWebpageRatePerSeconds int64
}

type SectionOpenAI struct {
	Secret    string
	Host      string
	ModelName string
}

type Config struct {
	TimezoneShiftSeconds int64
	Telegram             SectionTelegram
	OpenAI               SectionOpenAI
	Pinecone             SectionPinecone
	CloverDBPath         string
	DB                   SectionDB
	Slack                SectionSlack
	Discord              SectionDiscord
	Redis                SectionRedis
	LogLevel             string
	LogFilePath          string
	HardLimit            SectionHardLimit
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

		envLogLevel := getEnv(EnvLogLevel)

		redisDB, redisDBParseErr := strconv.ParseInt(getEnv(EnvRedisDB), 10, 64)
		log.Printf("failed to parse redis db %v: %v, should be number", getEnv(EnvRedisDB), redisDBParseErr)
		if redisDB < 0 {
			redisDB = 0

			log.Printf("redis db %v is less than 0, fallbacks to 0", getEnv(EnvRedisDB))
		}

		timezoneShiftSeconds, timezoneShiftSecondsParseErr := strconv.ParseInt(getEnv(EnvTimezoneShiftSeconds), 10, 64)
		log.Printf("failed to parse %s %v: %v, should be number", EnvTimezoneShiftSeconds, getEnv(EnvTimezoneShiftSeconds), timezoneShiftSecondsParseErr)

		manualRecapRatePerSecondsHardLimit, manualRecapHardLimitParseErr := strconv.ParseInt(getEnv(EnvHardLimitManualRecapRatePerSeconds), 10, 64)
		log.Printf("failed to parse %s %v: %v, should be number", EnvHardLimitManualRecapRatePerSeconds, getEnv(EnvHardLimitManualRecapRatePerSeconds), manualRecapHardLimitParseErr)

		if manualRecapRatePerSecondsHardLimit < 0 {
			manualRecapRatePerSecondsHardLimit = 0

			log.Printf("%s value %v is less than 0, fallbacks to 0", EnvHardLimitManualRecapRatePerSeconds, getEnv(EnvHardLimitManualRecapRatePerSeconds))
		}

		summarizeWebpageRatePerSecondsHardLimit, summarizeWebpageHardLimitParseErr := strconv.ParseInt(getEnv(EnvHardLimitSummarizeWebpageRatePerSeconds), 10, 64)
		log.Printf("failed to parse %s %v: %v, should be number", EnvHardLimitSummarizeWebpageRatePerSeconds, getEnv(EnvHardLimitSummarizeWebpageRatePerSeconds), summarizeWebpageHardLimitParseErr)

		if summarizeWebpageRatePerSecondsHardLimit < 0 {
			summarizeWebpageRatePerSecondsHardLimit = 0

			log.Printf("%s value %v is less than 0, fallbacks to 0", EnvHardLimitSummarizeWebpageRatePerSeconds, getEnv(EnvHardLimitSummarizeWebpageRatePerSeconds))
		}

		return &Config{
			TimezoneShiftSeconds: lo.Ternary(timezoneShiftSecondsParseErr == nil, lo.Ternary(timezoneShiftSeconds != 0, timezoneShiftSeconds, 0), 0),
			Telegram: SectionTelegram{
				BotToken:       getEnv(EnvTelegramBotToken),
				BotWebhookURL:  getEnv(EnvTelegramBotWebhookURL),
				BotWebhookPort: getEnv(EnvTelegramBotWebhookPort),
				BotAPIEndpoint: getEnv(EnvTelegramBotAPIEndpoint),
			},
			Slack: SectionSlack{
				Port:         getEnv(EnvSlackWebhookPort),
				ClientID:     getEnv(EnvSlackClientID),
				ClientSecret: getEnv(EnvSlackClientSecret),
			},
			OpenAI: SectionOpenAI{
				Secret:    getEnv(EnvOpenAIAPISecret),
				Host:      getEnv(EnvOpenAIAPIHost),
				ModelName: lo.Ternary(getEnv(EnvOpenAIAPIModelName) == "", goopenai.GPT3Dot5Turbo, getEnv(EnvOpenAIAPIModelName)),
			},
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
			Discord: SectionDiscord{
				Port:      getEnv(EnvDiscordBotWebhookPort),
				Token:     getEnv(EnvDiscordBotToken),
				PublicKey: getEnv(EnvDiscordBotPublicKey),
			},
			Redis: SectionRedis{
				Host:               getEnv(EnvRedisHost),
				Port:               getEnv(EnvRedisPort),
				TLSEnabled:         getEnv(EnvRedisTLSEnabled) == "true" || getEnv(EnvRedisTLSEnabled) == "1",
				Username:           getEnv(EnvRedisUsername),
				Password:           getEnv(EnvRedisPassword),
				DB:                 lo.Ternary(redisDBParseErr == nil, lo.Ternary(redisDB != 0, redisDB, 0), 0),
				ClientCacheEnabled: getEnv(EnvRedisClientCacheEnabled) == "true" || getEnv(EnvRedisClientCacheEnabled) == "1",
			},
			LogLevel:    lo.Ternary(envLogLevel == "", "info", envLogLevel),
			LogFilePath: getEnv(EnvLogFilePath),
			HardLimit: SectionHardLimit{
				ManualRecapRatePerSeconds:      manualRecapRatePerSecondsHardLimit,
				SummarizeWebpageRatePerSeconds: summarizeWebpageRatePerSecondsHardLimit,
			},
		}, nil
	}
}

func NewTestConfig() func() *Config {
	return func() *Config {
		return &Config{
			DB: SectionDB{
				ConnectionString: lo.Ternary(
					os.Getenv(EnvDBConnectionString) == "",
					"postgresql://postgres:123456@localhost:5432/postgres?search_path=public&sslmode=disable",
					os.Getenv(EnvDBConnectionString),
				),
			},
			Redis: SectionRedis{
				Host:               lo.Ternary(os.Getenv(EnvRedisHost) == "", "localhost", os.Getenv(EnvRedisHost)),
				Port:               lo.Ternary(os.Getenv(EnvRedisPort) == "", "6379", os.Getenv(EnvRedisPort)),
				TLSEnabled:         false,
				ClientCacheEnabled: false,
			},
			LogLevel: "debug",
		}
	}
}
