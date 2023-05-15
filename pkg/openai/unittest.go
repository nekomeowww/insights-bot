package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

var _ Client = (*MockClient)(nil)

type MockClient struct {
	SplitContentBasedByTokenLimitationsReturns []string

	TruncateContentBasedOnTokensReturns string

	SummarizeAnyReturns *openai.ChatCompletionResponse
	SummarizeAnyError   error

	SummarizeWithChatHistoriesReturns *openai.ChatCompletionResponse
	SummarizeWithOneChatHistoryError  error

	SummarizeWithOneChatHistoryReturns *openai.ChatCompletionResponse
	SummarizeWithChatHistoriesError    error

	SummarizeWithQuestionsAsSimplifiedChineseReturns *openai.ChatCompletionResponse
	SummarizeWithQuestionsAsSimplifiedChineseError   error
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m MockClient) SplitContentBasedByTokenLimitations(textContent string, limits int) []string {
	return m.SplitContentBasedByTokenLimitationsReturns
}

func (m MockClient) TruncateContentBasedOnTokens(textContent string, limits int) string {
	return m.TruncateContentBasedOnTokensReturns
}

func (m MockClient) SummarizeAny(ctx context.Context, content string) (*openai.ChatCompletionResponse, error) {
	if m.SummarizeAnyError != nil {
		return nil, m.SummarizeAnyError
	}

	return m.SummarizeAnyReturns, nil
}

func (m MockClient) SummarizeWithChatHistories(ctx context.Context, llmFriendlyChatHistories string) (*openai.ChatCompletionResponse, error) {
	if m.SummarizeWithChatHistoriesError != nil {
		return nil, m.SummarizeWithChatHistoriesError
	}

	return m.SummarizeWithChatHistoriesReturns, nil
}

func (m MockClient) SummarizeWithOneChatHistory(ctx context.Context, llmFriendlyChatHistory string) (*openai.ChatCompletionResponse, error) {
	if m.SummarizeWithOneChatHistoryError != nil {
		return nil, m.SummarizeWithOneChatHistoryError
	}

	return m.SummarizeWithOneChatHistoryReturns, nil
}

func (m MockClient) SummarizeWithQuestionsAsSimplifiedChinese(ctx context.Context, title string, by string, content string) (*openai.ChatCompletionResponse, error) {
	if m.SummarizeWithQuestionsAsSimplifiedChineseError != nil {
		return nil, m.SummarizeWithQuestionsAsSimplifiedChineseError
	}

	return m.SummarizeWithQuestionsAsSimplifiedChineseReturns, nil
}
