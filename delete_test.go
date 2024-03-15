package sq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteBuilderToSql(t *testing.T) {
	b := Delete("").
		Prefix("WITH prefix AS ?", 0).
		From("a").
		Where("b = ?", 1).
		OrderBy("c").
		Limit(2).
		Offset(3).
		Suffix("RETURNING ?", 4)

	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql :=
		"WITH prefix AS ? " +
			"DELETE FROM a WHERE b = ? ORDER BY c LIMIT 2 OFFSET 3 " +
			"RETURNING ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{0, 1, 4}
	require.Equal(t, expectedArgs, args)
}

func TestDeleteBuilderToSqlErr(t *testing.T) {
	_, _, err := Delete("").ToSql()
	require.Error(t, err)
}

func TestDeleteBuilderPlaceholders(t *testing.T) {
	b := Delete("test").Where("x = ? AND y = ?", 1, 2)

	sql, _, _ := b.PlaceholderFormat(Question).ToSql()
	require.Equal(t, "DELETE FROM test WHERE x = ? AND y = ?", sql)

	sql, _, _ = b.PlaceholderFormat(Dollar).ToSql()
	require.Equal(t, "DELETE FROM test WHERE x = $1 AND y = $2", sql)
}
