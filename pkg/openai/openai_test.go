package openai

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitContentBasedOnTokenLimitations(t *testing.T) {
	require := require.New(t)

	c := &Client{}
	slices, err := c.SplitContentBasedByTokenLimitations(strings.Repeat("a", 20000))
	require.NoError(err)

	require.Equal(3, len(slices))
}
