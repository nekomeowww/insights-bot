package summarize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractContentFromURL(t *testing.T) {
	extractContentFromURL("http://a.b.c")
}

func TestContentTypeCheck(t *testing.T) {
	assert.True(t, strings.Contains("text/html; charset=utf-8", "text/html"))
}
