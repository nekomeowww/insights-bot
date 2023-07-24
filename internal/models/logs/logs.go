package logs

import (
	"context"

	"github.com/nekomeowww/insights-bot/ent/logchathistoriesrecap"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewModelParams struct {
	fx.In

	Config *configs.Config
	Ent    *datastore.Ent
	Digger *datastore.AutoRecapTimeCapsuleDigger
	Logger *logger.Logger
}

type Model struct {
	config *configs.Config
	ent    *datastore.Ent
	logger *logger.Logger
	digger *datastore.AutoRecapTimeCapsuleDigger
}

func NewModel() func(NewModelParams) (*Model, error) {
	return func(param NewModelParams) (*Model, error) {
		return &Model{
			config: param.Config,
			ent:    param.Ent,
			logger: param.Logger,
			digger: param.Digger,
		}, nil
	}
}

func (m *Model) PruneAllLogsContentForChatID(chatID int64) error {
	err := m.ent.LogChatHistoriesRecap.
		Update().
		Where(
			logchathistoriesrecap.ChatIDEQ(chatID),
		).
		SetRecapInputs("").
		SetRecapOutputs("").
		Exec(context.TODO())

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) MigrateLogsOfChatFromChatIDToChatID(fromChatID int64, toChatID int64) error {
	affectedRows, err := m.ent.LogChatHistoriesRecap.
		Update().
		Where(
			logchathistoriesrecap.ChatIDEQ(fromChatID),
		).
		SetChatID(toChatID).
		Save(context.TODO())

	if err != nil {
		return err
	}

	m.logger.Info("successfully migrated options of chat",
		zap.Int64("from_chat_id", fromChatID),
		zap.Int64("to_chat_id", toChatID),
		zap.Int("affected_rows", affectedRows),
	)

	return nil
}
