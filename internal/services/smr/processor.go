package smr

import (
	"context"
	"errors"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
	"github.com/samber/lo"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func (s *Service) formatOutput(info types.TaskInfo, result *smr.URLSummarizationOutput) string {
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

func (s *Service) formatError(err error, language string) string {
	if errors.Is(err, smr.ErrContentNotSupported) {
		return s.i18n.TWithLanguage(language, "commands.groups.summarization.commands.smr.contentNotSupported")
	} else if errors.Is(err, smr.ErrNetworkError) || errors.Is(err, smr.ErrRequestFailed) {
		return s.i18n.TWithLanguage(language, "commands.groups.summarization.commands.smr.failedToReadDueToFailedToFetch")
	}

	return s.i18n.TWithLanguage(language, "commands.groups.summarization.commands.smr.failedToRead")
}

func (s *Service) newRetryButtonMarkup(info types.TaskInfo) (tgbotapi.InlineKeyboardMarkup, error) {
	data, err := s.tgBot.Bot().AssignOneCallbackQueryData("smr/summarization/retry", &info)

	if err != nil {
		return tgbotapi.InlineKeyboardMarkup{}, err
	}

	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		{
			Text:         s.i18n.TWithLanguage(info.Language, "commands.groups.summarization.commands.smr.retry"),
			CallbackData: &data,
		},
	}), nil
}

func (s *Service) sendResult(output *smr.URLSummarizationOutput, info types.TaskInfo, result string, provideRetryButton bool) {
	switch info.Platform {
	case bot.FromPlatformTelegram:
		msgEdit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    info.ChatID,
				MessageID: info.MessageID,
			},
			Text:      result,
			ParseMode: tgbotapi.ModeHTML,
		}

		if provideRetryButton {
			var err error
			retryButtonMarkup, err := s.newRetryButtonMarkup(info)

			if err != nil {
				s.logger.Error("smr service: failed to create retry button markup",
					zap.Error(err),
					zap.Int64("chat_id", info.ChatID),
					zap.String("platform", info.Platform.String()),
				)

				return
			}

			msgEdit.ReplyMarkup = &retryButtonMarkup
		}

		if output == nil {
			_, err := s.tgBot.Send(msgEdit)
			if err != nil {
				s.logger.Warn("smr service: failed to send result message",
					zap.Error(err),
					zap.Int64("chat_id", msgEdit.ChatID),
					zap.String("platform", info.Platform.String()),
				)
			}

			return
		}

		counts, err := s.model.FindFeedbackSummarizationsReactionCountsForChatIDAndLogID(info.ChatID, output.ID)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", info.ChatID),
				zap.String("platform", info.Platform.String()),
			)

			return
		}

		inlineKeyboardMarkup, err := s.model.NewVoteSummarizationsReactionsInlineKeyboardMarkup(s.tgBot.Bot(), info.ChatID, output.ID, counts.UpVotes, counts.DownVotes, counts.Lmao)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", info.ChatID),
				zap.String("platform", info.Platform.String()),
			)

			return
		}

		msgEdit.ReplyMarkup = lo.ToPtr(inlineKeyboardMarkup)

		_, err = s.tgBot.Send(msgEdit)
		if err != nil {
			s.logger.Warn("smr service: failed to send result message",
				zap.Error(err),
				zap.Int64("chat_id", msgEdit.ChatID),
				zap.String("platform", info.Platform.String()),
			)
		}
	case bot.FromPlatformSlack:
		// TODO: provide retry button
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
		// TODO: provide retry button
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
		errStr := s.formatError(err, lo.Ternary(info.Language == "", "en", info.Language))
		s.sendResult(nil, info, errStr, !errors.Is(err, smr.ErrContentNotSupported))

		return
	}

	finalResult := s.formatOutput(info, smrResult)
	s.sendResult(smrResult, info, finalResult, false)
}
