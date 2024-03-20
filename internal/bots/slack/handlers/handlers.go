package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"

	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot/services"
	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	types "github.com/nekomeowww/insights-bot/pkg/types/smr"
	"github.com/slack-go/slack"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModules() fx.Option {
	return fx.Options(
		fx.Provide(NewHandlers()),
	)
}

type NewHandlersParam struct {
	fx.In

	Config   *configs.Config
	Logger   *logger.Logger
	Ent      *datastore.Ent
	SmrQueue *smrqueue.Queue
	Services *services.Services
	I18n     *i18n.I18n
}

type Handlers struct {
	config   *configs.Config
	logger   *logger.Logger
	ent      *datastore.Ent
	smrQueue *smrqueue.Queue
	services *services.Services
	i18n     *i18n.I18n
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
			config:   param.Config,
			ent:      param.Ent,
			logger:   param.Logger,
			smrQueue: param.SmrQueue,
			services: param.Services,
			i18n:     param.I18n,
		}
	}
}

type receivedCommandInfo struct {
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseUrl string `form:"response_url"`
	UserID      string `form:"user_id"`
	ChannelID   string `form:"channel_id"`
	TeamID      string `form:"team_id"`
}

func (h *Handlers) PostCommandInfo(ctx *gin.Context) {
	var body receivedCommandInfo
	if err := ctx.Bind(&body); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Warn("failed to bind request body, type definition of slack request body may have changed", zap.Error(err))

		return
	}

	h.logger.Debug(fmt.Sprintf("slack: command received: /smr %s", body.Text),
		zap.String("user_id", body.UserID),
		zap.String("channel_id", body.ChannelID),
	)

	// get user locale, navie code, maybe need to refactor
	token, err := h.ent.SlackOAuthCredentials.Query().
		Where(slackoauthcredentials.TeamID(body.TeamID)).
		First(context.Background())

	if err != nil {
		h.logger.Warn("smr service: failed to get team's access token when get user locale",
			zap.Error(err),
		)

		return
	}

	slackCli := slackbot.NewSlackCli(nil, h.config.Slack.ClientID, h.config.Slack.ClientSecret, token.RefreshToken, token.AccessToken)
	user, err := slackCli.GetUserInfoWithTokenExpirationCheck(body.UserID, h.services.NewStoreFuncForRefresh(body.TeamID))

	if err != nil {
		h.logger.Warn("smr service: failed to user locale",
			zap.Error(err),
		)

		return
	}

	urlString := body.Text

	urlString = strings.TrimSpace(urlString)
	if !strings.HasPrefix(urlString, "http://") && !strings.HasPrefix(urlString, "https://") {
		urlString = "https://" + urlString
	}

	err, originErr := smr.CheckUrl(urlString)
	if err != nil {
		if smr.IsUrlCheckError(err) {
			ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(smr.FormatUrlCheckError(err, bot.FromPlatformSlack, user.Locale, nil)))
			return
		}

		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(h.i18n.TWithLanguage(user.Locale, "commands.groups.summarization.commands.smr.failedToRead")))
		h.logger.Warn("slack: failed to send error message", zap.Error(err), zap.NamedError("original_error", originErr))

		return
	}

	// check permissions
	_, err = h.ent.SlackOAuthCredentials.Query().Where(
		slackoauthcredentials.TeamID(body.TeamID),
	).First(context.Background())
	if err != nil {
		h.logger.Warn("slack: failed to get team's access token", zap.Error(err))
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(h.i18n.TWithLanguage(user.Locale, "commands.groups.summarization.commands.smr.permissionDenied")))
			return
		}

		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(h.i18n.TWithLanguage(user.Locale, "commands.groups.summarization.commands.smr.failedToRead")))

		return
	}

	// add task
	err = h.smrQueue.AddTask(types.TaskInfo{
		Platform:  bot.FromPlatformSlack,
		URL:       urlString,
		ChannelID: body.ChannelID,
		TeamID:    body.TeamID,
		Language:  user.Locale,
	})
	if err != nil {
		h.logger.Warn("slack: failed to add task", zap.Error(err))
		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(h.i18n.TWithLanguage(user.Locale, "commands.groups.summarization.commands.smr.failedToRead")))

		return
	}

	// response
	ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(h.i18n.TWithLanguage(user.Locale, "commands.groups.summarization.commands.smr.reading")))
}

// GetInstallAuth Receive auth code and request for access token.
func (h *Handlers) GetInstallAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := slack.GetOAuthV2Response(&http.Client{}, h.config.Slack.ClientID, h.config.Slack.ClientSecret, code, "")
	if err != nil {
		h.logger.Error("slack: failed to get access token, interrupt", zap.Error(err))
		ctx.AbortWithStatus(http.StatusServiceUnavailable)

		return
	}

	err = h.services.CreateOrUpdateSlackCredential(resp.Team.ID, resp.AccessToken, resp.RefreshToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("content-type", "text/html")
	_, _ = ctx.Writer.Write([]byte("<h1 style=\"text-align:center\">Success! Now you can close this page<h1>"))

	ctx.Status(http.StatusOK)
}
