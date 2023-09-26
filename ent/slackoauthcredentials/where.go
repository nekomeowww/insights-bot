// Code generated by ent, DO NOT EDIT.

package slackoauthcredentials

import (
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/nekomeowww/insights-bot/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldID, id))
}

// TeamID applies equality check predicate on the "team_id" field. It's identical to TeamIDEQ.
func TeamID(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldTeamID, v))
}

// RefreshToken applies equality check predicate on the "refresh_token" field. It's identical to RefreshTokenEQ.
func RefreshToken(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldRefreshToken, v))
}

// AccessToken applies equality check predicate on the "access_token" field. It's identical to AccessTokenEQ.
func AccessToken(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldAccessToken, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldUpdatedAt, v))
}

// TeamIDEQ applies the EQ predicate on the "team_id" field.
func TeamIDEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldTeamID, v))
}

// TeamIDNEQ applies the NEQ predicate on the "team_id" field.
func TeamIDNEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldTeamID, v))
}

// TeamIDIn applies the In predicate on the "team_id" field.
func TeamIDIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldTeamID, vs...))
}

// TeamIDNotIn applies the NotIn predicate on the "team_id" field.
func TeamIDNotIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldTeamID, vs...))
}

// TeamIDGT applies the GT predicate on the "team_id" field.
func TeamIDGT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldTeamID, v))
}

// TeamIDGTE applies the GTE predicate on the "team_id" field.
func TeamIDGTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldTeamID, v))
}

// TeamIDLT applies the LT predicate on the "team_id" field.
func TeamIDLT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldTeamID, v))
}

// TeamIDLTE applies the LTE predicate on the "team_id" field.
func TeamIDLTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldTeamID, v))
}

// TeamIDContains applies the Contains predicate on the "team_id" field.
func TeamIDContains(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContains(FieldTeamID, v))
}

// TeamIDHasPrefix applies the HasPrefix predicate on the "team_id" field.
func TeamIDHasPrefix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasPrefix(FieldTeamID, v))
}

// TeamIDHasSuffix applies the HasSuffix predicate on the "team_id" field.
func TeamIDHasSuffix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasSuffix(FieldTeamID, v))
}

// TeamIDEqualFold applies the EqualFold predicate on the "team_id" field.
func TeamIDEqualFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEqualFold(FieldTeamID, v))
}

// TeamIDContainsFold applies the ContainsFold predicate on the "team_id" field.
func TeamIDContainsFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContainsFold(FieldTeamID, v))
}

// RefreshTokenEQ applies the EQ predicate on the "refresh_token" field.
func RefreshTokenEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldRefreshToken, v))
}

// RefreshTokenNEQ applies the NEQ predicate on the "refresh_token" field.
func RefreshTokenNEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldRefreshToken, v))
}

// RefreshTokenIn applies the In predicate on the "refresh_token" field.
func RefreshTokenIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldRefreshToken, vs...))
}

// RefreshTokenNotIn applies the NotIn predicate on the "refresh_token" field.
func RefreshTokenNotIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldRefreshToken, vs...))
}

// RefreshTokenGT applies the GT predicate on the "refresh_token" field.
func RefreshTokenGT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldRefreshToken, v))
}

// RefreshTokenGTE applies the GTE predicate on the "refresh_token" field.
func RefreshTokenGTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldRefreshToken, v))
}

// RefreshTokenLT applies the LT predicate on the "refresh_token" field.
func RefreshTokenLT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldRefreshToken, v))
}

// RefreshTokenLTE applies the LTE predicate on the "refresh_token" field.
func RefreshTokenLTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldRefreshToken, v))
}

// RefreshTokenContains applies the Contains predicate on the "refresh_token" field.
func RefreshTokenContains(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContains(FieldRefreshToken, v))
}

// RefreshTokenHasPrefix applies the HasPrefix predicate on the "refresh_token" field.
func RefreshTokenHasPrefix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasPrefix(FieldRefreshToken, v))
}

// RefreshTokenHasSuffix applies the HasSuffix predicate on the "refresh_token" field.
func RefreshTokenHasSuffix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasSuffix(FieldRefreshToken, v))
}

// RefreshTokenEqualFold applies the EqualFold predicate on the "refresh_token" field.
func RefreshTokenEqualFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEqualFold(FieldRefreshToken, v))
}

// RefreshTokenContainsFold applies the ContainsFold predicate on the "refresh_token" field.
func RefreshTokenContainsFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContainsFold(FieldRefreshToken, v))
}

// AccessTokenEQ applies the EQ predicate on the "access_token" field.
func AccessTokenEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldAccessToken, v))
}

// AccessTokenNEQ applies the NEQ predicate on the "access_token" field.
func AccessTokenNEQ(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldAccessToken, v))
}

// AccessTokenIn applies the In predicate on the "access_token" field.
func AccessTokenIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldAccessToken, vs...))
}

// AccessTokenNotIn applies the NotIn predicate on the "access_token" field.
func AccessTokenNotIn(vs ...string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldAccessToken, vs...))
}

// AccessTokenGT applies the GT predicate on the "access_token" field.
func AccessTokenGT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldAccessToken, v))
}

// AccessTokenGTE applies the GTE predicate on the "access_token" field.
func AccessTokenGTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldAccessToken, v))
}

// AccessTokenLT applies the LT predicate on the "access_token" field.
func AccessTokenLT(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldAccessToken, v))
}

// AccessTokenLTE applies the LTE predicate on the "access_token" field.
func AccessTokenLTE(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldAccessToken, v))
}

// AccessTokenContains applies the Contains predicate on the "access_token" field.
func AccessTokenContains(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContains(FieldAccessToken, v))
}

// AccessTokenHasPrefix applies the HasPrefix predicate on the "access_token" field.
func AccessTokenHasPrefix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasPrefix(FieldAccessToken, v))
}

// AccessTokenHasSuffix applies the HasSuffix predicate on the "access_token" field.
func AccessTokenHasSuffix(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldHasSuffix(FieldAccessToken, v))
}

// AccessTokenEqualFold applies the EqualFold predicate on the "access_token" field.
func AccessTokenEqualFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEqualFold(FieldAccessToken, v))
}

// AccessTokenContainsFold applies the ContainsFold predicate on the "access_token" field.
func AccessTokenContainsFold(v string) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldContainsFold(FieldAccessToken, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v int64) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.FieldLTE(FieldUpdatedAt, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.SlackOAuthCredentials) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.SlackOAuthCredentials) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.SlackOAuthCredentials) predicate.SlackOAuthCredentials {
	return predicate.SlackOAuthCredentials(sql.NotPredicates(p))
}
