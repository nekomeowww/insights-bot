// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/telegramchatfeatureflags"
)

// TelegramChatFeatureFlagsCreate is the builder for creating a TelegramChatFeatureFlags entity.
type TelegramChatFeatureFlagsCreate struct {
	config
	mutation *TelegramChatFeatureFlagsMutation
	hooks    []Hook
}

// SetChatID sets the "chat_id" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetChatID(i int64) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetChatID(i)
	return tcffc
}

// SetChatType sets the "chat_type" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetChatType(s string) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetChatType(s)
	return tcffc
}

// SetFeatureChatHistoriesRecap sets the "feature_chat_histories_recap" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetFeatureChatHistoriesRecap(b bool) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetFeatureChatHistoriesRecap(b)
	return tcffc
}

// SetCreatedAt sets the "created_at" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetCreatedAt(i int64) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetCreatedAt(i)
	return tcffc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tcffc *TelegramChatFeatureFlagsCreate) SetNillableCreatedAt(i *int64) *TelegramChatFeatureFlagsCreate {
	if i != nil {
		tcffc.SetCreatedAt(*i)
	}
	return tcffc
}

// SetUpdatedAt sets the "updated_at" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetUpdatedAt(i int64) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetUpdatedAt(i)
	return tcffc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tcffc *TelegramChatFeatureFlagsCreate) SetNillableUpdatedAt(i *int64) *TelegramChatFeatureFlagsCreate {
	if i != nil {
		tcffc.SetUpdatedAt(*i)
	}
	return tcffc
}

// SetID sets the "id" field.
func (tcffc *TelegramChatFeatureFlagsCreate) SetID(u uuid.UUID) *TelegramChatFeatureFlagsCreate {
	tcffc.mutation.SetID(u)
	return tcffc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tcffc *TelegramChatFeatureFlagsCreate) SetNillableID(u *uuid.UUID) *TelegramChatFeatureFlagsCreate {
	if u != nil {
		tcffc.SetID(*u)
	}
	return tcffc
}

// Mutation returns the TelegramChatFeatureFlagsMutation object of the builder.
func (tcffc *TelegramChatFeatureFlagsCreate) Mutation() *TelegramChatFeatureFlagsMutation {
	return tcffc.mutation
}

// Save creates the TelegramChatFeatureFlags in the database.
func (tcffc *TelegramChatFeatureFlagsCreate) Save(ctx context.Context) (*TelegramChatFeatureFlags, error) {
	tcffc.defaults()
	return withHooks[*TelegramChatFeatureFlags, TelegramChatFeatureFlagsMutation](ctx, tcffc.sqlSave, tcffc.mutation, tcffc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tcffc *TelegramChatFeatureFlagsCreate) SaveX(ctx context.Context) *TelegramChatFeatureFlags {
	v, err := tcffc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcffc *TelegramChatFeatureFlagsCreate) Exec(ctx context.Context) error {
	_, err := tcffc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcffc *TelegramChatFeatureFlagsCreate) ExecX(ctx context.Context) {
	if err := tcffc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tcffc *TelegramChatFeatureFlagsCreate) defaults() {
	if _, ok := tcffc.mutation.CreatedAt(); !ok {
		v := telegramchatfeatureflags.DefaultCreatedAt()
		tcffc.mutation.SetCreatedAt(v)
	}
	if _, ok := tcffc.mutation.UpdatedAt(); !ok {
		v := telegramchatfeatureflags.DefaultUpdatedAt()
		tcffc.mutation.SetUpdatedAt(v)
	}
	if _, ok := tcffc.mutation.ID(); !ok {
		v := telegramchatfeatureflags.DefaultID()
		tcffc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcffc *TelegramChatFeatureFlagsCreate) check() error {
	if _, ok := tcffc.mutation.ChatID(); !ok {
		return &ValidationError{Name: "chat_id", err: errors.New(`ent: missing required field "TelegramChatFeatureFlags.chat_id"`)}
	}
	if _, ok := tcffc.mutation.ChatType(); !ok {
		return &ValidationError{Name: "chat_type", err: errors.New(`ent: missing required field "TelegramChatFeatureFlags.chat_type"`)}
	}
	if _, ok := tcffc.mutation.FeatureChatHistoriesRecap(); !ok {
		return &ValidationError{Name: "feature_chat_histories_recap", err: errors.New(`ent: missing required field "TelegramChatFeatureFlags.feature_chat_histories_recap"`)}
	}
	if _, ok := tcffc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "TelegramChatFeatureFlags.created_at"`)}
	}
	if _, ok := tcffc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "TelegramChatFeatureFlags.updated_at"`)}
	}
	return nil
}

func (tcffc *TelegramChatFeatureFlagsCreate) sqlSave(ctx context.Context) (*TelegramChatFeatureFlags, error) {
	if err := tcffc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tcffc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tcffc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	tcffc.mutation.id = &_node.ID
	tcffc.mutation.done = true
	return _node, nil
}

func (tcffc *TelegramChatFeatureFlagsCreate) createSpec() (*TelegramChatFeatureFlags, *sqlgraph.CreateSpec) {
	var (
		_node = &TelegramChatFeatureFlags{config: tcffc.config}
		_spec = sqlgraph.NewCreateSpec(telegramchatfeatureflags.Table, sqlgraph.NewFieldSpec(telegramchatfeatureflags.FieldID, field.TypeUUID))
	)
	_spec.Schema = tcffc.schemaConfig.TelegramChatFeatureFlags
	if id, ok := tcffc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tcffc.mutation.ChatID(); ok {
		_spec.SetField(telegramchatfeatureflags.FieldChatID, field.TypeInt64, value)
		_node.ChatID = value
	}
	if value, ok := tcffc.mutation.ChatType(); ok {
		_spec.SetField(telegramchatfeatureflags.FieldChatType, field.TypeString, value)
		_node.ChatType = value
	}
	if value, ok := tcffc.mutation.FeatureChatHistoriesRecap(); ok {
		_spec.SetField(telegramchatfeatureflags.FieldFeatureChatHistoriesRecap, field.TypeBool, value)
		_node.FeatureChatHistoriesRecap = value
	}
	if value, ok := tcffc.mutation.CreatedAt(); ok {
		_spec.SetField(telegramchatfeatureflags.FieldCreatedAt, field.TypeInt64, value)
		_node.CreatedAt = value
	}
	if value, ok := tcffc.mutation.UpdatedAt(); ok {
		_spec.SetField(telegramchatfeatureflags.FieldUpdatedAt, field.TypeInt64, value)
		_node.UpdatedAt = value
	}
	return _node, _spec
}

// TelegramChatFeatureFlagsCreateBulk is the builder for creating many TelegramChatFeatureFlags entities in bulk.
type TelegramChatFeatureFlagsCreateBulk struct {
	config
	builders []*TelegramChatFeatureFlagsCreate
}

// Save creates the TelegramChatFeatureFlags entities in the database.
func (tcffcb *TelegramChatFeatureFlagsCreateBulk) Save(ctx context.Context) ([]*TelegramChatFeatureFlags, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcffcb.builders))
	nodes := make([]*TelegramChatFeatureFlags, len(tcffcb.builders))
	mutators := make([]Mutator, len(tcffcb.builders))
	for i := range tcffcb.builders {
		func(i int, root context.Context) {
			builder := tcffcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TelegramChatFeatureFlagsMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tcffcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcffcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcffcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcffcb *TelegramChatFeatureFlagsCreateBulk) SaveX(ctx context.Context) []*TelegramChatFeatureFlags {
	v, err := tcffcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcffcb *TelegramChatFeatureFlagsCreateBulk) Exec(ctx context.Context) error {
	_, err := tcffcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcffcb *TelegramChatFeatureFlagsCreateBulk) ExecX(ctx context.Context) {
	if err := tcffcb.Exec(ctx); err != nil {
		panic(err)
	}
}
