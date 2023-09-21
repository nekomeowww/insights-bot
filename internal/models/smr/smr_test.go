package smr

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/configs"
	"github.com/nekomeowww/insights-bot/internal/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var model *Model

func TestMain(m *testing.M) {
	logger, err := lib.NewLogger()(lib.NewLoggerParams{
		Configs: configs.NewTestConfig()(),
	})
	if err != nil {
		panic(err)
	}

	model = NewModel()(NewModelParams{
		Logger: logger,
	})

	os.Exit(m.Run())
}

func TestExtractContentFromURL(t *testing.T) {
	t.Run("NoHost", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		article, err := model.extractContentFromURL(context.Background(), "https://a.b.c")
		require.Error(err)
		require.Nil(article)

		assert.ErrorIs(err, ErrNetworkError)
		assert.Contains(err.Error(), `no such host`)
	})

	t.Run("WeChatOfficialAccount", func(t *testing.T) {
		t.Skip("skip WeChatOfficialAccount, only test it when needed.")

		assert := assert.New(t)
		require := require.New(t)

		article, err := model.extractContentFromURL(context.Background(), fmt.Sprintf("https://mp.weixin.qq.com/s/%s", ""))
		require.NoError(err)

		assert.NotEmpty(article.Title)
		assert.NotEmpty(article.TextContent)
	})
}

func TestContentTypeCheck(t *testing.T) {
	assert.True(t, strings.Contains("text/html; charset=utf-8", "text/html"))
}
