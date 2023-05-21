// Code generated by counterfeiter. DO NOT EDIT.
package openaimock

import (
	"context"
	"sync"

	"github.com/nekomeowww/insights-bot/pkg/openai"
	openaia "github.com/sashabaranov/go-openai"
)

type MockClient struct {
	SplitContentBasedByTokenLimitationsStub        func(string, int) []string
	splitContentBasedByTokenLimitationsMutex       sync.RWMutex
	splitContentBasedByTokenLimitationsArgsForCall []struct {
		arg1 string
		arg2 int
	}
	splitContentBasedByTokenLimitationsReturns struct {
		result1 []string
	}
	splitContentBasedByTokenLimitationsReturnsOnCall map[int]struct {
		result1 []string
	}
	SummarizeAnyStub        func(context.Context, string) (*openaia.ChatCompletionResponse, error)
	summarizeAnyMutex       sync.RWMutex
	summarizeAnyArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	summarizeAnyReturns struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	summarizeAnyReturnsOnCall map[int]struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	SummarizeWithChatHistoriesStub        func(context.Context, string) (*openaia.ChatCompletionResponse, error)
	summarizeWithChatHistoriesMutex       sync.RWMutex
	summarizeWithChatHistoriesArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	summarizeWithChatHistoriesReturns struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	summarizeWithChatHistoriesReturnsOnCall map[int]struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	SummarizeWithOneChatHistoryStub        func(context.Context, string) (*openaia.ChatCompletionResponse, error)
	summarizeWithOneChatHistoryMutex       sync.RWMutex
	summarizeWithOneChatHistoryArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	summarizeWithOneChatHistoryReturns struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	summarizeWithOneChatHistoryReturnsOnCall map[int]struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	SummarizeWithQuestionsAsSimplifiedChineseStub        func(context.Context, string, string, string) (*openaia.ChatCompletionResponse, error)
	summarizeWithQuestionsAsSimplifiedChineseMutex       sync.RWMutex
	summarizeWithQuestionsAsSimplifiedChineseArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}
	summarizeWithQuestionsAsSimplifiedChineseReturns struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	summarizeWithQuestionsAsSimplifiedChineseReturnsOnCall map[int]struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}
	TruncateContentBasedOnTokensStub        func(string, int) string
	truncateContentBasedOnTokensMutex       sync.RWMutex
	truncateContentBasedOnTokensArgsForCall []struct {
		arg1 string
		arg2 int
	}
	truncateContentBasedOnTokensReturns struct {
		result1 string
	}
	truncateContentBasedOnTokensReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *MockClient) SplitContentBasedByTokenLimitations(arg1 string, arg2 int) []string {
	fake.splitContentBasedByTokenLimitationsMutex.Lock()
	ret, specificReturn := fake.splitContentBasedByTokenLimitationsReturnsOnCall[len(fake.splitContentBasedByTokenLimitationsArgsForCall)]
	fake.splitContentBasedByTokenLimitationsArgsForCall = append(fake.splitContentBasedByTokenLimitationsArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	stub := fake.SplitContentBasedByTokenLimitationsStub
	fakeReturns := fake.splitContentBasedByTokenLimitationsReturns
	fake.recordInvocation("SplitContentBasedByTokenLimitations", []interface{}{arg1, arg2})
	fake.splitContentBasedByTokenLimitationsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *MockClient) SplitContentBasedByTokenLimitationsCallCount() int {
	fake.splitContentBasedByTokenLimitationsMutex.RLock()
	defer fake.splitContentBasedByTokenLimitationsMutex.RUnlock()
	return len(fake.splitContentBasedByTokenLimitationsArgsForCall)
}

func (fake *MockClient) SplitContentBasedByTokenLimitationsCalls(stub func(string, int) []string) {
	fake.splitContentBasedByTokenLimitationsMutex.Lock()
	defer fake.splitContentBasedByTokenLimitationsMutex.Unlock()
	fake.SplitContentBasedByTokenLimitationsStub = stub
}

func (fake *MockClient) SplitContentBasedByTokenLimitationsArgsForCall(i int) (string, int) {
	fake.splitContentBasedByTokenLimitationsMutex.RLock()
	defer fake.splitContentBasedByTokenLimitationsMutex.RUnlock()
	argsForCall := fake.splitContentBasedByTokenLimitationsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MockClient) SplitContentBasedByTokenLimitationsReturns(result1 []string) {
	fake.splitContentBasedByTokenLimitationsMutex.Lock()
	defer fake.splitContentBasedByTokenLimitationsMutex.Unlock()
	fake.SplitContentBasedByTokenLimitationsStub = nil
	fake.splitContentBasedByTokenLimitationsReturns = struct {
		result1 []string
	}{result1}
}

func (fake *MockClient) SplitContentBasedByTokenLimitationsReturnsOnCall(i int, result1 []string) {
	fake.splitContentBasedByTokenLimitationsMutex.Lock()
	defer fake.splitContentBasedByTokenLimitationsMutex.Unlock()
	fake.SplitContentBasedByTokenLimitationsStub = nil
	if fake.splitContentBasedByTokenLimitationsReturnsOnCall == nil {
		fake.splitContentBasedByTokenLimitationsReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.splitContentBasedByTokenLimitationsReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *MockClient) SummarizeAny(arg1 context.Context, arg2 string) (*openaia.ChatCompletionResponse, error) {
	fake.summarizeAnyMutex.Lock()
	ret, specificReturn := fake.summarizeAnyReturnsOnCall[len(fake.summarizeAnyArgsForCall)]
	fake.summarizeAnyArgsForCall = append(fake.summarizeAnyArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.SummarizeAnyStub
	fakeReturns := fake.summarizeAnyReturns
	fake.recordInvocation("SummarizeAny", []interface{}{arg1, arg2})
	fake.summarizeAnyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockClient) SummarizeAnyCallCount() int {
	fake.summarizeAnyMutex.RLock()
	defer fake.summarizeAnyMutex.RUnlock()
	return len(fake.summarizeAnyArgsForCall)
}

func (fake *MockClient) SummarizeAnyCalls(stub func(context.Context, string) (*openaia.ChatCompletionResponse, error)) {
	fake.summarizeAnyMutex.Lock()
	defer fake.summarizeAnyMutex.Unlock()
	fake.SummarizeAnyStub = stub
}

func (fake *MockClient) SummarizeAnyArgsForCall(i int) (context.Context, string) {
	fake.summarizeAnyMutex.RLock()
	defer fake.summarizeAnyMutex.RUnlock()
	argsForCall := fake.summarizeAnyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MockClient) SummarizeAnyReturns(result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeAnyMutex.Lock()
	defer fake.summarizeAnyMutex.Unlock()
	fake.SummarizeAnyStub = nil
	fake.summarizeAnyReturns = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeAnyReturnsOnCall(i int, result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeAnyMutex.Lock()
	defer fake.summarizeAnyMutex.Unlock()
	fake.SummarizeAnyStub = nil
	if fake.summarizeAnyReturnsOnCall == nil {
		fake.summarizeAnyReturnsOnCall = make(map[int]struct {
			result1 *openaia.ChatCompletionResponse
			result2 error
		})
	}
	fake.summarizeAnyReturnsOnCall[i] = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithChatHistories(arg1 context.Context, arg2 string) (*openaia.ChatCompletionResponse, error) {
	fake.summarizeWithChatHistoriesMutex.Lock()
	ret, specificReturn := fake.summarizeWithChatHistoriesReturnsOnCall[len(fake.summarizeWithChatHistoriesArgsForCall)]
	fake.summarizeWithChatHistoriesArgsForCall = append(fake.summarizeWithChatHistoriesArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.SummarizeWithChatHistoriesStub
	fakeReturns := fake.summarizeWithChatHistoriesReturns
	fake.recordInvocation("SummarizeWithChatHistories", []interface{}{arg1, arg2})
	fake.summarizeWithChatHistoriesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockClient) SummarizeWithChatHistoriesCallCount() int {
	fake.summarizeWithChatHistoriesMutex.RLock()
	defer fake.summarizeWithChatHistoriesMutex.RUnlock()
	return len(fake.summarizeWithChatHistoriesArgsForCall)
}

func (fake *MockClient) SummarizeWithChatHistoriesCalls(stub func(context.Context, string) (*openaia.ChatCompletionResponse, error)) {
	fake.summarizeWithChatHistoriesMutex.Lock()
	defer fake.summarizeWithChatHistoriesMutex.Unlock()
	fake.SummarizeWithChatHistoriesStub = stub
}

func (fake *MockClient) SummarizeWithChatHistoriesArgsForCall(i int) (context.Context, string) {
	fake.summarizeWithChatHistoriesMutex.RLock()
	defer fake.summarizeWithChatHistoriesMutex.RUnlock()
	argsForCall := fake.summarizeWithChatHistoriesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MockClient) SummarizeWithChatHistoriesReturns(result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithChatHistoriesMutex.Lock()
	defer fake.summarizeWithChatHistoriesMutex.Unlock()
	fake.SummarizeWithChatHistoriesStub = nil
	fake.summarizeWithChatHistoriesReturns = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithChatHistoriesReturnsOnCall(i int, result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithChatHistoriesMutex.Lock()
	defer fake.summarizeWithChatHistoriesMutex.Unlock()
	fake.SummarizeWithChatHistoriesStub = nil
	if fake.summarizeWithChatHistoriesReturnsOnCall == nil {
		fake.summarizeWithChatHistoriesReturnsOnCall = make(map[int]struct {
			result1 *openaia.ChatCompletionResponse
			result2 error
		})
	}
	fake.summarizeWithChatHistoriesReturnsOnCall[i] = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithOneChatHistory(arg1 context.Context, arg2 string) (*openaia.ChatCompletionResponse, error) {
	fake.summarizeWithOneChatHistoryMutex.Lock()
	ret, specificReturn := fake.summarizeWithOneChatHistoryReturnsOnCall[len(fake.summarizeWithOneChatHistoryArgsForCall)]
	fake.summarizeWithOneChatHistoryArgsForCall = append(fake.summarizeWithOneChatHistoryArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.SummarizeWithOneChatHistoryStub
	fakeReturns := fake.summarizeWithOneChatHistoryReturns
	fake.recordInvocation("SummarizeWithOneChatHistory", []interface{}{arg1, arg2})
	fake.summarizeWithOneChatHistoryMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockClient) SummarizeWithOneChatHistoryCallCount() int {
	fake.summarizeWithOneChatHistoryMutex.RLock()
	defer fake.summarizeWithOneChatHistoryMutex.RUnlock()
	return len(fake.summarizeWithOneChatHistoryArgsForCall)
}

func (fake *MockClient) SummarizeWithOneChatHistoryCalls(stub func(context.Context, string) (*openaia.ChatCompletionResponse, error)) {
	fake.summarizeWithOneChatHistoryMutex.Lock()
	defer fake.summarizeWithOneChatHistoryMutex.Unlock()
	fake.SummarizeWithOneChatHistoryStub = stub
}

func (fake *MockClient) SummarizeWithOneChatHistoryArgsForCall(i int) (context.Context, string) {
	fake.summarizeWithOneChatHistoryMutex.RLock()
	defer fake.summarizeWithOneChatHistoryMutex.RUnlock()
	argsForCall := fake.summarizeWithOneChatHistoryArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MockClient) SummarizeWithOneChatHistoryReturns(result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithOneChatHistoryMutex.Lock()
	defer fake.summarizeWithOneChatHistoryMutex.Unlock()
	fake.SummarizeWithOneChatHistoryStub = nil
	fake.summarizeWithOneChatHistoryReturns = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithOneChatHistoryReturnsOnCall(i int, result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithOneChatHistoryMutex.Lock()
	defer fake.summarizeWithOneChatHistoryMutex.Unlock()
	fake.SummarizeWithOneChatHistoryStub = nil
	if fake.summarizeWithOneChatHistoryReturnsOnCall == nil {
		fake.summarizeWithOneChatHistoryReturnsOnCall = make(map[int]struct {
			result1 *openaia.ChatCompletionResponse
			result2 error
		})
	}
	fake.summarizeWithOneChatHistoryReturnsOnCall[i] = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChinese(arg1 context.Context, arg2 string, arg3 string, arg4 string) (*openaia.ChatCompletionResponse, error) {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Lock()
	ret, specificReturn := fake.summarizeWithQuestionsAsSimplifiedChineseReturnsOnCall[len(fake.summarizeWithQuestionsAsSimplifiedChineseArgsForCall)]
	fake.summarizeWithQuestionsAsSimplifiedChineseArgsForCall = append(fake.summarizeWithQuestionsAsSimplifiedChineseArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	stub := fake.SummarizeWithQuestionsAsSimplifiedChineseStub
	fakeReturns := fake.summarizeWithQuestionsAsSimplifiedChineseReturns
	fake.recordInvocation("SummarizeWithQuestionsAsSimplifiedChinese", []interface{}{arg1, arg2, arg3, arg4})
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChineseCallCount() int {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RLock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RUnlock()
	return len(fake.summarizeWithQuestionsAsSimplifiedChineseArgsForCall)
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChineseCalls(stub func(context.Context, string, string, string) (*openaia.ChatCompletionResponse, error)) {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Lock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Unlock()
	fake.SummarizeWithQuestionsAsSimplifiedChineseStub = stub
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChineseArgsForCall(i int) (context.Context, string, string, string) {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RLock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RUnlock()
	argsForCall := fake.summarizeWithQuestionsAsSimplifiedChineseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChineseReturns(result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Lock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Unlock()
	fake.SummarizeWithQuestionsAsSimplifiedChineseStub = nil
	fake.summarizeWithQuestionsAsSimplifiedChineseReturns = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) SummarizeWithQuestionsAsSimplifiedChineseReturnsOnCall(i int, result1 *openaia.ChatCompletionResponse, result2 error) {
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Lock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.Unlock()
	fake.SummarizeWithQuestionsAsSimplifiedChineseStub = nil
	if fake.summarizeWithQuestionsAsSimplifiedChineseReturnsOnCall == nil {
		fake.summarizeWithQuestionsAsSimplifiedChineseReturnsOnCall = make(map[int]struct {
			result1 *openaia.ChatCompletionResponse
			result2 error
		})
	}
	fake.summarizeWithQuestionsAsSimplifiedChineseReturnsOnCall[i] = struct {
		result1 *openaia.ChatCompletionResponse
		result2 error
	}{result1, result2}
}

func (fake *MockClient) TruncateContentBasedOnTokens(arg1 string, arg2 int) string {
	fake.truncateContentBasedOnTokensMutex.Lock()
	ret, specificReturn := fake.truncateContentBasedOnTokensReturnsOnCall[len(fake.truncateContentBasedOnTokensArgsForCall)]
	fake.truncateContentBasedOnTokensArgsForCall = append(fake.truncateContentBasedOnTokensArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	stub := fake.TruncateContentBasedOnTokensStub
	fakeReturns := fake.truncateContentBasedOnTokensReturns
	fake.recordInvocation("TruncateContentBasedOnTokens", []interface{}{arg1, arg2})
	fake.truncateContentBasedOnTokensMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *MockClient) TruncateContentBasedOnTokensCallCount() int {
	fake.truncateContentBasedOnTokensMutex.RLock()
	defer fake.truncateContentBasedOnTokensMutex.RUnlock()
	return len(fake.truncateContentBasedOnTokensArgsForCall)
}

func (fake *MockClient) TruncateContentBasedOnTokensCalls(stub func(string, int) string) {
	fake.truncateContentBasedOnTokensMutex.Lock()
	defer fake.truncateContentBasedOnTokensMutex.Unlock()
	fake.TruncateContentBasedOnTokensStub = stub
}

func (fake *MockClient) TruncateContentBasedOnTokensArgsForCall(i int) (string, int) {
	fake.truncateContentBasedOnTokensMutex.RLock()
	defer fake.truncateContentBasedOnTokensMutex.RUnlock()
	argsForCall := fake.truncateContentBasedOnTokensArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *MockClient) TruncateContentBasedOnTokensReturns(result1 string) {
	fake.truncateContentBasedOnTokensMutex.Lock()
	defer fake.truncateContentBasedOnTokensMutex.Unlock()
	fake.TruncateContentBasedOnTokensStub = nil
	fake.truncateContentBasedOnTokensReturns = struct {
		result1 string
	}{result1}
}

func (fake *MockClient) TruncateContentBasedOnTokensReturnsOnCall(i int, result1 string) {
	fake.truncateContentBasedOnTokensMutex.Lock()
	defer fake.truncateContentBasedOnTokensMutex.Unlock()
	fake.TruncateContentBasedOnTokensStub = nil
	if fake.truncateContentBasedOnTokensReturnsOnCall == nil {
		fake.truncateContentBasedOnTokensReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.truncateContentBasedOnTokensReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *MockClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.splitContentBasedByTokenLimitationsMutex.RLock()
	defer fake.splitContentBasedByTokenLimitationsMutex.RUnlock()
	fake.summarizeAnyMutex.RLock()
	defer fake.summarizeAnyMutex.RUnlock()
	fake.summarizeWithChatHistoriesMutex.RLock()
	defer fake.summarizeWithChatHistoriesMutex.RUnlock()
	fake.summarizeWithOneChatHistoryMutex.RLock()
	defer fake.summarizeWithOneChatHistoryMutex.RUnlock()
	fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RLock()
	defer fake.summarizeWithQuestionsAsSimplifiedChineseMutex.RUnlock()
	fake.truncateContentBasedOnTokensMutex.RLock()
	defer fake.truncateContentBasedOnTokensMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *MockClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ openai.Client = new(MockClient)
