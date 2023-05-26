package chathistories

import (
	"strings"
	"testing"

	"github.com/nekomeowww/insights-bot/internal/thirdparty/openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecapOutputTemplateExecute(t *testing.T) {
	sb := new(strings.Builder)
	err := RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
		ChatID: formatChatID(-100123456789),
		Recap: &openai.ChatHistorySummarizationOutputs{
			TopicName:                        "Topic 1",
			SinceID:                          1,
			ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
			Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
				{
					Point:  "Point 1",
					KeyIDs: []int64{1, 2},
				},
				{
					Point: "Point 2",
				},
			},
			Conclusion: "Conclusion 1",
		},
	})
	require.NoError(t, err)
	expected := `## <a href="https://t.me/c/123456789/1">Topic 1</a>
参与人：User 1，User 2
讨论：
 - Point 1 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>
 - Point 2
结论：Conclusion 1`
	assert.Equal(t, expected, sb.String())

	sb = new(strings.Builder)
	err = RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
		ChatID: formatChatID(-100123456789),
		Recap: &openai.ChatHistorySummarizationOutputs{
			TopicName:                        "Topic 3",
			ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
			Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
				{
					Point: "Point 1",
				},
				{
					Point:  "Point 2",
					KeyIDs: []int64{1, 2},
				},
			},
		},
	})
	require.NoError(t, err)
	expected = `## Topic 3
参与人：User 1，User 2
讨论：
 - Point 1
 - Point 2 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>`
	assert.Equal(t, expected, sb.String())

	sb = new(strings.Builder)
	err = RecapOutputTemplate.Execute(sb, RecapOutputTemplateInputs{
		ChatID: formatChatID(-100123456789),
		Recap: &openai.ChatHistorySummarizationOutputs{
			TopicName:                        "Topic 1",
			SinceID:                          2,
			ParticipantsNamesWithoutUsername: []string{"User 1", "User 2"},
			Discussion: []*openai.ChatHistorySummarizationOutputsDiscussion{
				{
					Point:  "Point 1",
					KeyIDs: []int64{1, 2},
				},
				{
					Point: "Point 2",
				},
			},
			Conclusion: "Conclusion 2",
		},
	})
	require.NoError(t, err)

	expected = `## <a href="https://t.me/c/123456789/2">Topic 1</a>
参与人：User 1，User 2
讨论：
 - Point 1 <a href="https://t.me/c/123456789/1">[1]</a> <a href="https://t.me/c/123456789/2">[2]</a>
 - Point 2
结论：Conclusion 2`
	assert.Equal(t, expected, sb.String())
}

func TestFormatFullNameAndUsername(t *testing.T) {
	tests := []struct {
		name     string
		fullName string
		username string
		result   string
	}{
		{
			name:     `full name shorter than 10 chars`,
			fullName: "Full Name",
			username: "example_username",
			result:   "Full Name",
		},
		{
			name:     `full name longer than 10 chars`,
			fullName: "A Very Long Full Name",
			username: "example_username",
			result:   "example_username",
		},
		{
			name:     `full name longer than 10 chars AND username is empty`,
			fullName: "A Very Long Full Name",
			username: "",
			result:   "A Very Long Full Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFullNameAndUsername(tt.fullName, tt.username)
			assert.Equal(t, tt.result, result)
		})
	}
}
