package smr

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
	"github.com/samber/lo"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (s *Service) processOutput(info types.TaskInfo, result *smr.URLSummarizationOutput) string {
	switch info.Platform {
	case bot.FromPlatformTelegram:
		return result.FormatSummarizationAsHTML()
	case bot.FromPlatformSlack:
		return result.FormatSummarizationAsSlackMarkdown()
	case bot.FromPlatformDiscord:
		return result.FormatSummarizationAsDiscordMarkdown()
	default:
		return ""
	}
}

func (s *Service) processError(err error) string {
	if errors.Is(err, smr.ErrContentNotSupported) {
		return "暂时不支持量子速读这样的内容呢，可以换个别的链接试试。"
	} else if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
		return "量子速读的链接读取失败了哦。可以再试试？"
	}

	return "量子速读失败了。可以再试试？"
}

func (s *Service) sendResult(output *smr.URLSummarizationOutput, info types.TaskInfo, result string) {
	switch info.Platform {
	case bot.FromPlatformTelegram:
		logID := uuid.Nil
		if output != nil {
			logID = output.ID
		}

		counts, err := s.model.FindFeedbackSummarizationsReactionCountsForChatIDAndLogID(info.ChatID, logID)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", info.ChatID),
				zap.String("platform", info.Platform.String()),
			)

			return
		}

		inlineKeyboardMarkup, err := s.model.NewVoteSummarizationsReactionsInlineKeyboardMarkup(s.tgBot.Bot(), info.ChatID, logID, counts.UpVotes, counts.DownVotes, counts.Lmao)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", info.ChatID),
				zap.String("platform", info.Platform.String()),
			)

			return
		}

		msgEdit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      info.ChatID,
				MessageID:   info.MessageID,
				ReplyMarkup: lo.ToPtr(inlineKeyboardMarkup),
			},
			Text:      result,
			ParseMode: tgbotapi.ModeHTML,
		}

		_, err = s.tgBot.Send(msgEdit)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", msgEdit.ChatID),
				zap.String("platform", info.Platform.String()),
			)
		}
	case bot.FromPlatformSlack:
		token, err := s.ent.SlackOAuthCredentials.Query().
			Where(slackoauthcredentials.TeamID(info.TeamID)).
			First(context.Background())

		if err != nil {
			s.logger.Warn("smr service: failed to get team's access token",
				zap.Error(err),
				zap.String("platform", info.Platform.String()),
			)

			return
		}

		slackCfg := s.config.Slack
		slackCli := slackbot.NewSlackCli(nil, slackCfg.ClientID, slackCfg.ClientSecret, token.RefreshToken, token.AccessToken)
		_, _, _, err = slackCli.SendMessageWithTokenExpirationCheck(
			info.ChannelID,
			s.slackBot.GetService().NewStoreFuncForRefresh(info.TeamID),
			slack.MsgOptionText(result, false),
		)

		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.String("platform", info.Platform.String()),
			)
		}
	case bot.FromPlatformDiscord:
		channelID, _ := snowflake.Parse(info.ChannelID)
		_, err := s.discordBot.Rest().
			CreateMessage(channelID, discord.NewMessageCreateBuilder().
				SetContent(result).
				Build(),
			)

		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.String("platform", info.Platform.String()),
			)
		}
	}
}

func (s *Service) isBotExists(platform bot.FromPlatform) bool {
	switch platform {
	case bot.FromPlatformTelegram:
		return s.tgBot != nil
	case bot.FromPlatformSlack:
		return s.slackBot != nil
	case bot.FromPlatformDiscord:
		return s.discordBot != nil
	}

	return false
}

func (s *Service) processor(info types.TaskInfo) {
	if !s.isBotExists(info.Platform) {
		s.logger.Error("received task from platform " + info.Platform.String() + " but instance not exists")
		// move back to queue
		err := s.queue.AddTask(info)
		if err != nil {
			s.logger.Error("failed to move task back to queue", zap.Error(err))
		}

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	smrResult, err := s.model.SummarizeInputURL(ctx, info.URL, info.Platform)
	if err != nil {
		s.logger.Warn("smr service: summarization failed", zap.Error(err))
		errStr := s.processError(err)
		s.sendResult(nil, info, errStr)

		return
	}

	finalResult := s.processOutput(info, smrResult)
	s.sendResult(smrResult, info, finalResult)
}
