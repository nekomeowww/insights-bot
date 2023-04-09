package openai

import (
	"strconv"
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
			textContent: "小溪河水清澈见底",
			limits:      3,
			expected:    "小溪",
		},
		{
			textContent: "小溪河水清澈见底",
			limits:      4,
			expected:    "小溪",
		},
		{
			textContent: "小溪河水清澈见底",
			limits:      5,
			expected:    "小溪河",
		},
	}

	c, err := NewClient("", "")
	require.NoError(t, err)
	for _, table := range tables {
		t.Run(table.textContent, func(t *testing.T) {
			actual := c.TruncateContentBasedOnTokens(table.textContent, table.limits)
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
			textContent: strings.Repeat("a", 40000),
			limits:      3900,
			expected:    []string{strings.Repeat("a", 31200), strings.Repeat("a", 8800)},
		},
		{
			textContent: "小溪河水清澈见底，沿岸芦苇丛生。远处山峰耸立，白云飘渺。一只黄鹂停在枝头，唱起了优美的歌曲，引来了不少路人驻足欣赏。",
			limits:      20,
			expected:    []string{"小溪河水清澈见底，沿岸芦", "苇丛生。远处山峰耸立，白", "云飘渺。一只黄鹂停在枝头，", "唱起了优美的歌曲，引来了不少路人", "驻足欣赏。"},
		},
	}

	c, err := NewClient("", "")
	require.NoError(t, err)
	for i, table := range tables {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := c.SplitContentBasedByTokenLimitations(table.textContent, table.limits)
			require.Equal(t, table.expected, actual)
		})
	}
}
