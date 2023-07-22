// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nekomeowww/insights-bot/ent/internal"
	"github.com/nekomeowww/insights-bot/ent/logsummarizations"
	"github.com/nekomeowww/insights-bot/ent/predicate"
)

// LogSummarizationsUpdate is the builder for updating LogSummarizations entities.
type LogSummarizationsUpdate struct {
	config
	hooks    []Hook
	mutation *LogSummarizationsMutation
}

// Where appends a list predicates to the LogSummarizationsUpdate builder.
func (lsu *LogSummarizationsUpdate) Where(ps ...predicate.LogSummarizations) *LogSummarizationsUpdate {
	lsu.mutation.Where(ps...)
	return lsu
}

// SetContentURL sets the "content_url" field.
func (lsu *LogSummarizationsUpdate) SetContentURL(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetContentURL(s)
	return lsu
}

// SetNillableContentURL sets the "content_url" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableContentURL(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetContentURL(*s)
	}
	return lsu
}

// SetContentTitle sets the "content_title" field.
func (lsu *LogSummarizationsUpdate) SetContentTitle(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetContentTitle(s)
	return lsu
}

// SetNillableContentTitle sets the "content_title" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableContentTitle(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetContentTitle(*s)
	}
	return lsu
}

// SetContentAuthor sets the "content_author" field.
func (lsu *LogSummarizationsUpdate) SetContentAuthor(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetContentAuthor(s)
	return lsu
}

// SetNillableContentAuthor sets the "content_author" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableContentAuthor(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetContentAuthor(*s)
	}
	return lsu
}

// SetContentText sets the "content_text" field.
func (lsu *LogSummarizationsUpdate) SetContentText(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetContentText(s)
	return lsu
}

// SetNillableContentText sets the "content_text" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableContentText(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetContentText(*s)
	}
	return lsu
}

// SetContentSummarizedOutputs sets the "content_summarized_outputs" field.
func (lsu *LogSummarizationsUpdate) SetContentSummarizedOutputs(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetContentSummarizedOutputs(s)
	return lsu
}

// SetNillableContentSummarizedOutputs sets the "content_summarized_outputs" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableContentSummarizedOutputs(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetContentSummarizedOutputs(*s)
	}
	return lsu
}

// SetFromPlatform sets the "from_platform" field.
func (lsu *LogSummarizationsUpdate) SetFromPlatform(i int) *LogSummarizationsUpdate {
	lsu.mutation.ResetFromPlatform()
	lsu.mutation.SetFromPlatform(i)
	return lsu
}

// SetNillableFromPlatform sets the "from_platform" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableFromPlatform(i *int) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetFromPlatform(*i)
	}
	return lsu
}

// AddFromPlatform adds i to the "from_platform" field.
func (lsu *LogSummarizationsUpdate) AddFromPlatform(i int) *LogSummarizationsUpdate {
	lsu.mutation.AddFromPlatform(i)
	return lsu
}

// SetPromptTokenUsage sets the "prompt_token_usage" field.
func (lsu *LogSummarizationsUpdate) SetPromptTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.ResetPromptTokenUsage()
	lsu.mutation.SetPromptTokenUsage(i)
	return lsu
}

// SetNillablePromptTokenUsage sets the "prompt_token_usage" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillablePromptTokenUsage(i *int) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetPromptTokenUsage(*i)
	}
	return lsu
}

// AddPromptTokenUsage adds i to the "prompt_token_usage" field.
func (lsu *LogSummarizationsUpdate) AddPromptTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.AddPromptTokenUsage(i)
	return lsu
}

// SetCompletionTokenUsage sets the "completion_token_usage" field.
func (lsu *LogSummarizationsUpdate) SetCompletionTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.ResetCompletionTokenUsage()
	lsu.mutation.SetCompletionTokenUsage(i)
	return lsu
}

// SetNillableCompletionTokenUsage sets the "completion_token_usage" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableCompletionTokenUsage(i *int) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetCompletionTokenUsage(*i)
	}
	return lsu
}

// AddCompletionTokenUsage adds i to the "completion_token_usage" field.
func (lsu *LogSummarizationsUpdate) AddCompletionTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.AddCompletionTokenUsage(i)
	return lsu
}

// SetTotalTokenUsage sets the "total_token_usage" field.
func (lsu *LogSummarizationsUpdate) SetTotalTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.ResetTotalTokenUsage()
	lsu.mutation.SetTotalTokenUsage(i)
	return lsu
}

// SetNillableTotalTokenUsage sets the "total_token_usage" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableTotalTokenUsage(i *int) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetTotalTokenUsage(*i)
	}
	return lsu
}

// AddTotalTokenUsage adds i to the "total_token_usage" field.
func (lsu *LogSummarizationsUpdate) AddTotalTokenUsage(i int) *LogSummarizationsUpdate {
	lsu.mutation.AddTotalTokenUsage(i)
	return lsu
}

// SetModelName sets the "model_name" field.
func (lsu *LogSummarizationsUpdate) SetModelName(s string) *LogSummarizationsUpdate {
	lsu.mutation.SetModelName(s)
	return lsu
}

// SetNillableModelName sets the "model_name" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableModelName(s *string) *LogSummarizationsUpdate {
	if s != nil {
		lsu.SetModelName(*s)
	}
	return lsu
}

// SetCreatedAt sets the "created_at" field.
func (lsu *LogSummarizationsUpdate) SetCreatedAt(i int64) *LogSummarizationsUpdate {
	lsu.mutation.ResetCreatedAt()
	lsu.mutation.SetCreatedAt(i)
	return lsu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableCreatedAt(i *int64) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetCreatedAt(*i)
	}
	return lsu
}

// AddCreatedAt adds i to the "created_at" field.
func (lsu *LogSummarizationsUpdate) AddCreatedAt(i int64) *LogSummarizationsUpdate {
	lsu.mutation.AddCreatedAt(i)
	return lsu
}

// SetUpdatedAt sets the "updated_at" field.
func (lsu *LogSummarizationsUpdate) SetUpdatedAt(i int64) *LogSummarizationsUpdate {
	lsu.mutation.ResetUpdatedAt()
	lsu.mutation.SetUpdatedAt(i)
	return lsu
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (lsu *LogSummarizationsUpdate) SetNillableUpdatedAt(i *int64) *LogSummarizationsUpdate {
	if i != nil {
		lsu.SetUpdatedAt(*i)
	}
	return lsu
}

// AddUpdatedAt adds i to the "updated_at" field.
func (lsu *LogSummarizationsUpdate) AddUpdatedAt(i int64) *LogSummarizationsUpdate {
	lsu.mutation.AddUpdatedAt(i)
	return lsu
}

// Mutation returns the LogSummarizationsMutation object of the builder.
func (lsu *LogSummarizationsUpdate) Mutation() *LogSummarizationsMutation {
	return lsu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (lsu *LogSummarizationsUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, lsu.sqlSave, lsu.mutation, lsu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lsu *LogSummarizationsUpdate) SaveX(ctx context.Context) int {
	affected, err := lsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (lsu *LogSummarizationsUpdate) Exec(ctx context.Context) error {
	_, err := lsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lsu *LogSummarizationsUpdate) ExecX(ctx context.Context) {
	if err := lsu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lsu *LogSummarizationsUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(logsummarizations.Table, logsummarizations.Columns, sqlgraph.NewFieldSpec(logsummarizations.FieldID, field.TypeUUID))
	if ps := lsu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lsu.mutation.ContentURL(); ok {
		_spec.SetField(logsummarizations.FieldContentURL, field.TypeString, value)
	}
	if value, ok := lsu.mutation.ContentTitle(); ok {
		_spec.SetField(logsummarizations.FieldContentTitle, field.TypeString, value)
	}
	if value, ok := lsu.mutation.ContentAuthor(); ok {
		_spec.SetField(logsummarizations.FieldContentAuthor, field.TypeString, value)
	}
	if value, ok := lsu.mutation.ContentText(); ok {
		_spec.SetField(logsummarizations.FieldContentText, field.TypeString, value)
	}
	if value, ok := lsu.mutation.ContentSummarizedOutputs(); ok {
		_spec.SetField(logsummarizations.FieldContentSummarizedOutputs, field.TypeString, value)
	}
	if value, ok := lsu.mutation.FromPlatform(); ok {
		_spec.SetField(logsummarizations.FieldFromPlatform, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.AddedFromPlatform(); ok {
		_spec.AddField(logsummarizations.FieldFromPlatform, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.PromptTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldPromptTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.AddedPromptTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldPromptTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.CompletionTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldCompletionTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.AddedCompletionTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldCompletionTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.TotalTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldTotalTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.AddedTotalTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldTotalTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsu.mutation.ModelName(); ok {
		_spec.SetField(logsummarizations.FieldModelName, field.TypeString, value)
	}
	if value, ok := lsu.mutation.CreatedAt(); ok {
		_spec.SetField(logsummarizations.FieldCreatedAt, field.TypeInt64, value)
	}
	if value, ok := lsu.mutation.AddedCreatedAt(); ok {
		_spec.AddField(logsummarizations.FieldCreatedAt, field.TypeInt64, value)
	}
	if value, ok := lsu.mutation.UpdatedAt(); ok {
		_spec.SetField(logsummarizations.FieldUpdatedAt, field.TypeInt64, value)
	}
	if value, ok := lsu.mutation.AddedUpdatedAt(); ok {
		_spec.AddField(logsummarizations.FieldUpdatedAt, field.TypeInt64, value)
	}
	_spec.Node.Schema = lsu.schemaConfig.LogSummarizations
	ctx = internal.NewSchemaConfigContext(ctx, lsu.schemaConfig)
	if n, err = sqlgraph.UpdateNodes(ctx, lsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{logsummarizations.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	lsu.mutation.done = true
	return n, nil
}

// LogSummarizationsUpdateOne is the builder for updating a single LogSummarizations entity.
type LogSummarizationsUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *LogSummarizationsMutation
}

// SetContentURL sets the "content_url" field.
func (lsuo *LogSummarizationsUpdateOne) SetContentURL(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetContentURL(s)
	return lsuo
}

// SetNillableContentURL sets the "content_url" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableContentURL(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetContentURL(*s)
	}
	return lsuo
}

// SetContentTitle sets the "content_title" field.
func (lsuo *LogSummarizationsUpdateOne) SetContentTitle(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetContentTitle(s)
	return lsuo
}

// SetNillableContentTitle sets the "content_title" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableContentTitle(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetContentTitle(*s)
	}
	return lsuo
}

// SetContentAuthor sets the "content_author" field.
func (lsuo *LogSummarizationsUpdateOne) SetContentAuthor(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetContentAuthor(s)
	return lsuo
}

// SetNillableContentAuthor sets the "content_author" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableContentAuthor(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetContentAuthor(*s)
	}
	return lsuo
}

// SetContentText sets the "content_text" field.
func (lsuo *LogSummarizationsUpdateOne) SetContentText(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetContentText(s)
	return lsuo
}

// SetNillableContentText sets the "content_text" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableContentText(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetContentText(*s)
	}
	return lsuo
}

// SetContentSummarizedOutputs sets the "content_summarized_outputs" field.
func (lsuo *LogSummarizationsUpdateOne) SetContentSummarizedOutputs(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetContentSummarizedOutputs(s)
	return lsuo
}

// SetNillableContentSummarizedOutputs sets the "content_summarized_outputs" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableContentSummarizedOutputs(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetContentSummarizedOutputs(*s)
	}
	return lsuo
}

// SetFromPlatform sets the "from_platform" field.
func (lsuo *LogSummarizationsUpdateOne) SetFromPlatform(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetFromPlatform()
	lsuo.mutation.SetFromPlatform(i)
	return lsuo
}

// SetNillableFromPlatform sets the "from_platform" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableFromPlatform(i *int) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetFromPlatform(*i)
	}
	return lsuo
}

// AddFromPlatform adds i to the "from_platform" field.
func (lsuo *LogSummarizationsUpdateOne) AddFromPlatform(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddFromPlatform(i)
	return lsuo
}

// SetPromptTokenUsage sets the "prompt_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) SetPromptTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetPromptTokenUsage()
	lsuo.mutation.SetPromptTokenUsage(i)
	return lsuo
}

// SetNillablePromptTokenUsage sets the "prompt_token_usage" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillablePromptTokenUsage(i *int) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetPromptTokenUsage(*i)
	}
	return lsuo
}

// AddPromptTokenUsage adds i to the "prompt_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) AddPromptTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddPromptTokenUsage(i)
	return lsuo
}

// SetCompletionTokenUsage sets the "completion_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) SetCompletionTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetCompletionTokenUsage()
	lsuo.mutation.SetCompletionTokenUsage(i)
	return lsuo
}

// SetNillableCompletionTokenUsage sets the "completion_token_usage" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableCompletionTokenUsage(i *int) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetCompletionTokenUsage(*i)
	}
	return lsuo
}

// AddCompletionTokenUsage adds i to the "completion_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) AddCompletionTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddCompletionTokenUsage(i)
	return lsuo
}

// SetTotalTokenUsage sets the "total_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) SetTotalTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetTotalTokenUsage()
	lsuo.mutation.SetTotalTokenUsage(i)
	return lsuo
}

// SetNillableTotalTokenUsage sets the "total_token_usage" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableTotalTokenUsage(i *int) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetTotalTokenUsage(*i)
	}
	return lsuo
}

// AddTotalTokenUsage adds i to the "total_token_usage" field.
func (lsuo *LogSummarizationsUpdateOne) AddTotalTokenUsage(i int) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddTotalTokenUsage(i)
	return lsuo
}

// SetModelName sets the "model_name" field.
func (lsuo *LogSummarizationsUpdateOne) SetModelName(s string) *LogSummarizationsUpdateOne {
	lsuo.mutation.SetModelName(s)
	return lsuo
}

// SetNillableModelName sets the "model_name" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableModelName(s *string) *LogSummarizationsUpdateOne {
	if s != nil {
		lsuo.SetModelName(*s)
	}
	return lsuo
}

// SetCreatedAt sets the "created_at" field.
func (lsuo *LogSummarizationsUpdateOne) SetCreatedAt(i int64) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetCreatedAt()
	lsuo.mutation.SetCreatedAt(i)
	return lsuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableCreatedAt(i *int64) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetCreatedAt(*i)
	}
	return lsuo
}

// AddCreatedAt adds i to the "created_at" field.
func (lsuo *LogSummarizationsUpdateOne) AddCreatedAt(i int64) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddCreatedAt(i)
	return lsuo
}

// SetUpdatedAt sets the "updated_at" field.
func (lsuo *LogSummarizationsUpdateOne) SetUpdatedAt(i int64) *LogSummarizationsUpdateOne {
	lsuo.mutation.ResetUpdatedAt()
	lsuo.mutation.SetUpdatedAt(i)
	return lsuo
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (lsuo *LogSummarizationsUpdateOne) SetNillableUpdatedAt(i *int64) *LogSummarizationsUpdateOne {
	if i != nil {
		lsuo.SetUpdatedAt(*i)
	}
	return lsuo
}

// AddUpdatedAt adds i to the "updated_at" field.
func (lsuo *LogSummarizationsUpdateOne) AddUpdatedAt(i int64) *LogSummarizationsUpdateOne {
	lsuo.mutation.AddUpdatedAt(i)
	return lsuo
}

// Mutation returns the LogSummarizationsMutation object of the builder.
func (lsuo *LogSummarizationsUpdateOne) Mutation() *LogSummarizationsMutation {
	return lsuo.mutation
}

// Where appends a list predicates to the LogSummarizationsUpdate builder.
func (lsuo *LogSummarizationsUpdateOne) Where(ps ...predicate.LogSummarizations) *LogSummarizationsUpdateOne {
	lsuo.mutation.Where(ps...)
	return lsuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (lsuo *LogSummarizationsUpdateOne) Select(field string, fields ...string) *LogSummarizationsUpdateOne {
	lsuo.fields = append([]string{field}, fields...)
	return lsuo
}

// Save executes the query and returns the updated LogSummarizations entity.
func (lsuo *LogSummarizationsUpdateOne) Save(ctx context.Context) (*LogSummarizations, error) {
	return withHooks(ctx, lsuo.sqlSave, lsuo.mutation, lsuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (lsuo *LogSummarizationsUpdateOne) SaveX(ctx context.Context) *LogSummarizations {
	node, err := lsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (lsuo *LogSummarizationsUpdateOne) Exec(ctx context.Context) error {
	_, err := lsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lsuo *LogSummarizationsUpdateOne) ExecX(ctx context.Context) {
	if err := lsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (lsuo *LogSummarizationsUpdateOne) sqlSave(ctx context.Context) (_node *LogSummarizations, err error) {
	_spec := sqlgraph.NewUpdateSpec(logsummarizations.Table, logsummarizations.Columns, sqlgraph.NewFieldSpec(logsummarizations.FieldID, field.TypeUUID))
	id, ok := lsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "LogSummarizations.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := lsuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, logsummarizations.FieldID)
		for _, f := range fields {
			if !logsummarizations.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != logsummarizations.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := lsuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := lsuo.mutation.ContentURL(); ok {
		_spec.SetField(logsummarizations.FieldContentURL, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.ContentTitle(); ok {
		_spec.SetField(logsummarizations.FieldContentTitle, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.ContentAuthor(); ok {
		_spec.SetField(logsummarizations.FieldContentAuthor, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.ContentText(); ok {
		_spec.SetField(logsummarizations.FieldContentText, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.ContentSummarizedOutputs(); ok {
		_spec.SetField(logsummarizations.FieldContentSummarizedOutputs, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.FromPlatform(); ok {
		_spec.SetField(logsummarizations.FieldFromPlatform, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.AddedFromPlatform(); ok {
		_spec.AddField(logsummarizations.FieldFromPlatform, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.PromptTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldPromptTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.AddedPromptTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldPromptTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.CompletionTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldCompletionTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.AddedCompletionTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldCompletionTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.TotalTokenUsage(); ok {
		_spec.SetField(logsummarizations.FieldTotalTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.AddedTotalTokenUsage(); ok {
		_spec.AddField(logsummarizations.FieldTotalTokenUsage, field.TypeInt, value)
	}
	if value, ok := lsuo.mutation.ModelName(); ok {
		_spec.SetField(logsummarizations.FieldModelName, field.TypeString, value)
	}
	if value, ok := lsuo.mutation.CreatedAt(); ok {
		_spec.SetField(logsummarizations.FieldCreatedAt, field.TypeInt64, value)
	}
	if value, ok := lsuo.mutation.AddedCreatedAt(); ok {
		_spec.AddField(logsummarizations.FieldCreatedAt, field.TypeInt64, value)
	}
	if value, ok := lsuo.mutation.UpdatedAt(); ok {
		_spec.SetField(logsummarizations.FieldUpdatedAt, field.TypeInt64, value)
	}
	if value, ok := lsuo.mutation.AddedUpdatedAt(); ok {
		_spec.AddField(logsummarizations.FieldUpdatedAt, field.TypeInt64, value)
	}
	_spec.Node.Schema = lsuo.schemaConfig.LogSummarizations
	ctx = internal.NewSchemaConfigContext(ctx, lsuo.schemaConfig)
	_node = &LogSummarizations{config: lsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, lsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{logsummarizations.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	lsuo.mutation.done = true
	return _node, nil
}
