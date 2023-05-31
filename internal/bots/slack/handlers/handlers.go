package handlers

import (
	"context"
	"github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrqueue"
	"github.com/nekomeowww/insights-bot/internal/services/smr/smrutils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nekomeowww/insights-bot/ent"
	"github.com/nekomeowww/insights-bot/ent/slackoauthcredentials"
	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/datastore"
	"github.com/nekomeowww/insights-bot/internal/services/smr/types"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot"
	"github.com/nekomeowww/insights-bot/pkg/bots/slackbot/services"
	"github.com/nekomeowww/insights-bot/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"go.uber.org/fx"
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
}

type Handlers struct {
	config   *configs.Config
	logger   *logger.Logger
	ent      *datastore.Ent
	smrQueue *smrqueue.Queue
	services *services.Services
}

func NewHandlers() func(param NewHandlersParam) *Handlers {
	return func(param NewHandlersParam) *Handlers {
		return &Handlers{
			config:   param.Config,
			ent:      param.Ent,
			logger:   param.Logger,
			smrQueue: param.SmrQueue,
			services: param.Services,
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
		h.logger.WithField("error", err.Error()).Warn("failed to bind request body, maybe slack request definition changed")

		return
	}

	h.logger.WithFields(logrus.Fields{
		"user_id":    body.UserID,
		"channel_id": body.ChannelID,
	}).Infof("slack: command received: /smr %s", body.Text)

	urlString := body.Text

	err := smrutils.CheckUrl(urlString)
	if err != nil {
		if smrutils.IsUrlCheckError(err) {
			ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage(err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("出现了一些问题，可以再试试？"))

		return
	}

	// check permissions
	_, err = h.ent.SlackOAuthCredentials.Query().Where(
		slackoauthcredentials.TeamID(body.TeamID),
	).First(context.Background())
	if err != nil {
		h.logger.WithField("error", err.Error()).Warn("slack: failed to get team's access token")
		if ent.IsNotFound(err) {
			ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("本应用没有权限向这个频道发送消息，尝试重新安装一下？"))
			return
		}

		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("出现了一些问题，可以再试试？"))

		return
	}

	// add task
	err = h.smrQueue.AddTask(types.TaskInfo{
		Platform:  smr.FromPlatformSlack,
		Url:       urlString,
		ChannelID: body.ChannelID,
		TeamID:    body.TeamID,
	})
	if err != nil {
		h.logger.WithError(err).Warn("slack: failed to add task")
		ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("量子速读请求发送失败了，可以再试试？"))

		return
	}

	// response
	ctx.JSON(http.StatusOK, slackbot.NewSlackWebhookMessage("请稍等，量子速读中..."))
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
		h.logger.WithError(err).Error("slack: failed to get access token, interrupt")
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
