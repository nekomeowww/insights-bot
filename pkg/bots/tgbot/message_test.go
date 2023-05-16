package tgbot

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitMessagesAgainstLengthLimitIntoMessageGroups(t *testing.T) {
	s := []string{
		strings.Repeat("a", 1005),
		strings.Repeat("b", 1005),
		strings.Repeat("c", 1005),
		strings.Repeat("d", 1005),
		strings.Repeat("e", 1005),
		strings.Repeat("f", 1005),
		strings.Repeat("g", 1005),
		strings.Repeat("h", 1005),
	}

	batchSlice := SplitMessagesAgainstLengthLimitIntoMessageGroups(s)
	require.Len(t, batchSlice, 3)

	require.Len(t, batchSlice[0], 3)
	require.Len(t, batchSlice[1], 3)
	require.Len(t, batchSlice[2], 2)
}
