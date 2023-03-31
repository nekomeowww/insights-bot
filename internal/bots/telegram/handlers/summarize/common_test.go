package summarize

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractContentFromURL(t *testing.T) {
	t.Run("NoHost", func(t *testing.T) {
		assert := assert.New(t)
		require := require.New(t)

		article, err := extractContentFromURL("https://a.b.c")
		require.Error(err)
		require.Nil(article)

		assert.ErrorIs(err, ErrNetworkError)
		assert.EqualError(err, `failed to get url https://a.b.c, network error: Get "https://a.b.c": dial tcp: lookup a.b.c: no such host`)
	})

	t.Run("WeChatOfficialAccount", func(t *testing.T) {
		t.Skip("skip WeChatOfficialAccount, only test it when needed.")

		assert := assert.New(t)
		require := require.New(t)

		article, err := extractContentFromURL(fmt.Sprintf("https://mp.weixin.qq.com/s/%s", ""))
		require.NoError(err)

		assert.NotEmpty(article.Title)
		assert.NotEmpty(article.TextContent)
	})
}

func TestContentTypeCheck(t *testing.T) {
	assert.True(t, strings.Contains("text/html; charset=utf-8", "text/html"))
}
