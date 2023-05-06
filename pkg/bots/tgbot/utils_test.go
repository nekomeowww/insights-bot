package tgbot

import (
	"fmt"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

func TestReplaceMarkdownTitlesToBoldTexts(t *testing.T) {
	prefix := ""
	for i := 0; i < 6; i++ {
		t.Run(fmt.Sprintf("TitleLevel%d", i+1), func(t *testing.T) {
			a := assert.New(t)

			prefix += "#"
			expect := "<b>Test</b>"
			actual, err := ReplaceMarkdownTitlesToTelegramBoldElement(fmt.Sprintf("%s Test", prefix))
			a.Nil(err)
			a.Equal(expect, actual)
		})
	}

	t.Run("MoreHashTags", func(t *testing.T) {
		a := assert.New(t)

		prefix += "#"
		expect := "####### Test"
		actual, err := ReplaceMarkdownTitlesToTelegramBoldElement(fmt.Sprintf("%s Test", prefix))
		a.Nil(err)
		a.Equal(expect, actual)
	})

	t.Run("MultipleLines", func(t *testing.T) {
		a := assert.New(t)

		expect := `<b>there is a title</b>
<b>there is a subtitle</b>`
		actual, err := ReplaceMarkdownTitlesToTelegramBoldElement(`# there is a title
## there is a subtitle`)
		a.Nil(err)
		a.Equal(expect, actual)
	})
}

func TestExtractTextFromMessage(t *testing.T) {
	t.Run("MixedUrlsAndTextLinks", func(t *testing.T) {
		a := assert.New(t)
		message := &tgbotapi.Message{
			MessageID: 666,
			From:      &tgbotapi.User{ID: 23333333},
			Date:      1683386000,
			Chat:      &tgbotapi.Chat{ID: 0xc0001145e4},
			Text:      "看看这些链接：https://docs.swift.org/swift-book/documentation/the-swift-programming-language/stringsandcharacters/#Extended-Grapheme-Clusters 、https://www.youtube.com/watch?v=outcGtbnMuQ https://github.com/nekomeowww/insights-bot 还有 这个",
			Entities: []tgbotapi.MessageEntity{
				{Type: "url", Offset: 7, Length: 127, URL: "", Language: ""},
				{Type: "url", Offset: 136, Length: 43, URL: "", Language: ""},
				{Type: "url", Offset: 180, Length: 42, URL: "", Language: ""},
				{Type: "text_link", Offset: 226, Length: 2, URL: "https://matters.town/@1435Club/322889-%E8%BF%99%E5%87%A0%E5%A4%A9-web3%E5%9C%A8%E5%A4%A7%E7%90%86%E5%8F%91%E7%94%9F%E4%BA%86%E4%BB%80%E4%B9%88", Language: ""},
			},
			Photo: []tgbotapi.PhotoSize{},
		}
		expect := "看看这些链接：[Documentation](https://docs.swift.org/swift-book/documentation/the-swift-programming-language/stringsandcharacters/#Extended-Grapheme-Clusters) 、[GPT-4 Developer Livestream](https://www.youtube.com/watch?v=outcGtbnMuQ) [GitHub - nekomeowww/insights-bot: A bot works with OpenAI GPT models to provide insights for your info flows.](https://github.com/nekomeowww/insights-bot) 还有 [这个](https://matters.town/@1435Club/322889-这几天-web3在大理发生了什么)"
		a.Equal(expect, ExtractTextFromMessage(message))
	})
}
