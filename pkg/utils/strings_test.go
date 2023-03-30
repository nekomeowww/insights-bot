package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceMarkdownTitlesToBoldTexts(t *testing.T) {
	prefix := ""
	for i := 0; i < 6; i++ {
		t.Run(fmt.Sprintf("title level %d", i+1), func(t *testing.T) {
			a := assert.New(t)

			prefix += "#"
			expect := "**Test**"
			actual, err := ReplaceMarkdownTitlesToBoldTexts(fmt.Sprintf("%s Test", prefix))
			a.Nil(err)
			a.Equal(expect, actual)
		})
	}

	t.Run("more hash tags", func(t *testing.T) {
		a := assert.New(t)

		prefix += "#"
		expect := "####### Test"
		actual, err := ReplaceMarkdownTitlesToBoldTexts(fmt.Sprintf("%s Test", prefix))
		a.Nil(err)
		a.Equal(expect, actual)
	})

	t.Run("multiple lines", func(t *testing.T) {
		a := assert.New(t)

		expect := `**there is a title**
**there is a subtitle**`
		actual, err := ReplaceMarkdownTitlesToBoldTexts(`# there is a title
## there is a subtitle`)
		a.Nil(err)
		a.Equal(expect, actual)
	})
}
