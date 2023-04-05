package openai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	NewClient("", "https://openai.example.com")
}

func TestSplitContentBasedOnTokenLimitations(t *testing.T) {
	require := require.New(t)

	c := &Client{}
	slices, err := c.SplitContentBasedByTokenLimitations(strings.Repeat("a", 20000), 3900)
	require.NoError(err)

	require.Equal(3, len(slices))
}
