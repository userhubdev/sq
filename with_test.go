package sq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithBuilder_ToSql(t *testing.T) {
	withClause := With().
		As("c", Select("a").From("b").Where("a = ?", 1))

	b := Select("a").PrefixExpr(withClause).From("c")
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH c AS ( SELECT a FROM b WHERE a = ?) SELECT a FROM c", sql)
	require.Equal(t, []any{1}, args)
}

func TestWithBuilder_SelectNoCTE(t *testing.T) {
	b := With().Select("a").From("b").Where("a = ?", 1)
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "SELECT a FROM b WHERE a = ?", sql)
	require.Equal(t, []any{1}, args)
}

func TestWithBuilder_SelectOneCTE(t *testing.T) {
	b := StatementBuilder.PlaceholderFormat(AtP).With().
		As("c", Select("a").From("b").Where("a = ?", 1)).
		Select("a").From("c")
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH c AS ( SELECT a FROM b WHERE a = @p1) SELECT a FROM c", sql)
	require.Equal(t, []any{1}, args)
}

func TestWithBuilder_SelectMultipleCTE(t *testing.T) {
	b := StatementBuilder.PlaceholderFormat(AtP).With().
		As("c", Select("a").From("b").Where("a > ?", 1)).
		As("d", Select("a").From("c").Where("a < ?", 100)).
		Select("a").From("d")
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH "+
		"c AS ( SELECT a FROM b WHERE a > @p1), "+
		"d AS ( SELECT a FROM c WHERE a < @p2) "+
		"SELECT a FROM d", sql)
	require.Equal(t, []any{1, 100}, args)
}

func TestWithBuilder_InsertNoCTE(t *testing.T) {
	b := With().Insert("b").SetMap(map[string]any{"a": 5})
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "INSERT INTO b (a) VALUES (?)", sql)
	require.Equal(t, []any{5}, args)
}

func TestWithBuilder_InsertOneCTE(t *testing.T) {
	b := StatementBuilder.PlaceholderFormat(AtP).With().
		As("c", Select("a").From("b").Where("a = ?", 1)).
		Insert("d").Columns("v").
		Select(Select("a").From("c"))
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH c AS ( SELECT a FROM b WHERE a = @p1) INSERT INTO d (v) SELECT a FROM c", sql)
	require.Equal(t, []any{1}, args)
}

func TestWithBuilder_UpdateNoCTE(t *testing.T) {
	b := With().Update("b").SetMap(map[string]any{"a": 5})
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "UPDATE b SET a = ?", sql)
	require.Equal(t, []any{5}, args)
}

func TestWithBuilder_UpdateOneCTE(t *testing.T) {
	b := StatementBuilder.PlaceholderFormat(AtP).With().
		As("c", Select("a").From("b").Where("a = ?", 1)).
		Update("d").
		SetMap(map[string]any{"v": 5}).
		Where(ConcatExpr("a IN (", Select("a").From("c"), ")"))
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH c AS ( SELECT a FROM b WHERE a = @p1) UPDATE d SET v = @p2 WHERE a IN (SELECT a FROM c)", sql)
	require.Equal(t, []any{1, 5}, args)
}

func TestWithBuilder_DeleteNoCTE(t *testing.T) {
	b := With().Delete("b").Where(Eq{"a": 5})
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "DELETE FROM b WHERE a = ?", sql)
	require.Equal(t, []any{5}, args)
}

func TestWithBuilder_DeleteOneCTE(t *testing.T) {
	b := StatementBuilder.PlaceholderFormat(AtP).With().
		As("c", Select("a").From("b").Where("a = ?", 1)).
		Delete("d").
		Where(ConcatExpr("a IN (", Select("a").From("c"), ")"))
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Equal(t, "WITH c AS ( SELECT a FROM b WHERE a = @p1) DELETE FROM d WHERE a IN (SELECT a FROM c)", sql)
	require.Equal(t, []any{1}, args)
}
