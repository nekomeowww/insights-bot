package tgbot

import (
	"fmt"
	"testing"

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
