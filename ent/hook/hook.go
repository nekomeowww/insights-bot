// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/nekomeowww/insights-bot/ent"
)

// The ChatHistoriesFunc type is an adapter to allow the use of ordinary
// function as ChatHistories mutator.
type ChatHistoriesFunc func(context.Context, *ent.ChatHistoriesMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ChatHistoriesFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ChatHistoriesMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ChatHistoriesMutation", m)
}

// The FeedbackChatHistoriesRecapsReactionsFunc type is an adapter to allow the use of ordinary
// function as FeedbackChatHistoriesRecapsReactions mutator.
type FeedbackChatHistoriesRecapsReactionsFunc func(context.Context, *ent.FeedbackChatHistoriesRecapsReactionsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FeedbackChatHistoriesRecapsReactionsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.FeedbackChatHistoriesRecapsReactionsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FeedbackChatHistoriesRecapsReactionsMutation", m)
}

// The FeedbackSummarizationsReactionsFunc type is an adapter to allow the use of ordinary
// function as FeedbackSummarizationsReactions mutator.
type FeedbackSummarizationsReactionsFunc func(context.Context, *ent.FeedbackSummarizationsReactionsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FeedbackSummarizationsReactionsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.FeedbackSummarizationsReactionsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FeedbackSummarizationsReactionsMutation", m)
}

// The LogChatHistoriesRecapFunc type is an adapter to allow the use of ordinary
// function as LogChatHistoriesRecap mutator.
type LogChatHistoriesRecapFunc func(context.Context, *ent.LogChatHistoriesRecapMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LogChatHistoriesRecapFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.LogChatHistoriesRecapMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LogChatHistoriesRecapMutation", m)
}

// The LogSummarizationsFunc type is an adapter to allow the use of ordinary
// function as LogSummarizations mutator.
type LogSummarizationsFunc func(context.Context, *ent.LogSummarizationsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LogSummarizationsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.LogSummarizationsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LogSummarizationsMutation", m)
}

// The MetricOpenAIChatCompletionTokenUsageFunc type is an adapter to allow the use of ordinary
// function as MetricOpenAIChatCompletionTokenUsage mutator.
type MetricOpenAIChatCompletionTokenUsageFunc func(context.Context, *ent.MetricOpenAIChatCompletionTokenUsageMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f MetricOpenAIChatCompletionTokenUsageFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.MetricOpenAIChatCompletionTokenUsageMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.MetricOpenAIChatCompletionTokenUsageMutation", m)
}

// The SentMessagesFunc type is an adapter to allow the use of ordinary
// function as SentMessages mutator.
type SentMessagesFunc func(context.Context, *ent.SentMessagesMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SentMessagesFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.SentMessagesMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SentMessagesMutation", m)
}

// The SlackOAuthCredentialsFunc type is an adapter to allow the use of ordinary
// function as SlackOAuthCredentials mutator.
type SlackOAuthCredentialsFunc func(context.Context, *ent.SlackOAuthCredentialsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SlackOAuthCredentialsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.SlackOAuthCredentialsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SlackOAuthCredentialsMutation", m)
}

// The TelegramChatAutoRecapsSubscribersFunc type is an adapter to allow the use of ordinary
// function as TelegramChatAutoRecapsSubscribers mutator.
type TelegramChatAutoRecapsSubscribersFunc func(context.Context, *ent.TelegramChatAutoRecapsSubscribersMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TelegramChatAutoRecapsSubscribersFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.TelegramChatAutoRecapsSubscribersMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TelegramChatAutoRecapsSubscribersMutation", m)
}

// The TelegramChatFeatureFlagsFunc type is an adapter to allow the use of ordinary
// function as TelegramChatFeatureFlags mutator.
type TelegramChatFeatureFlagsFunc func(context.Context, *ent.TelegramChatFeatureFlagsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TelegramChatFeatureFlagsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.TelegramChatFeatureFlagsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TelegramChatFeatureFlagsMutation", m)
}

// The TelegramChatRecapsOptionsFunc type is an adapter to allow the use of ordinary
// function as TelegramChatRecapsOptions mutator.
type TelegramChatRecapsOptionsFunc func(context.Context, *ent.TelegramChatRecapsOptionsMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TelegramChatRecapsOptionsFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.TelegramChatRecapsOptionsMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TelegramChatRecapsOptionsMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
