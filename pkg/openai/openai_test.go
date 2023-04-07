package openai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	NewClient("", "https://openai.example.com")
}

func TestTruncateContentBasedOnTokens(t *testing.T) {
	tables := []struct {
		textContent string
		limits      int
		expected    string
	}{
		{
			textContent: "心理学家",
			limits:      4,
			expected:    "心理",
		},
		{
			textContent: "心理学家",
			limits:      5,
			expected:    "心理",
		},
		{
			textContent: "心理学家",
			limits:      6,
			expected:    "心理学",
		},
		{
			textContent: "心理学家",
			limits:      10,
			expected:    "心理学家",
		},
	}

	for _, table := range tables {
		t.Run(table.textContent, func(t *testing.T) {
			c := &Client{}
			actual, err := c.TruncateContentBasedOnTokens(table.textContent, table.limits)
			require.NoError(t, err)
			require.Equal(t, table.expected, actual)
		})
	}
}

func TestSplitContentBasedOnTokenLimitations(t *testing.T) {
	tables := []struct {
		textContent string
		limits      int
		expected    []string
	}{
		{
			textContent: strings.Repeat("a", 20000),
			limits:      3900,
			expected:    []string{strings.Repeat("a", 15600), strings.Repeat("a", 4400)},
		},
		{
			textContent: "小溪河水清澈见底，沿岸芦苇丛生。远处山峰耸立，白云飘渺。一只黄鹂停在枝头，唱起了优美的歌曲，引来了不少路人驻足欣赏。",
			limits:      20,
			expected:    []string{"小溪河水清澈见底", "，沿岸芦苇丛生。", "远处山峰耸立，白", "云飘渺。一只黄鹂停", "在枝头，唱起了优", "美的歌曲，引来了不", "少路人驻足欣赏。"},
		},
	}

	for _, table := range tables {
		t.Run(table.textContent, func(t *testing.T) {
			c := &Client{}
			actual, err := c.SplitContentBasedByTokenLimitations(table.textContent, table.limits)
			require.NoError(t, err)
			require.Equal(t, table.expected, actual)
		})
	}
}
