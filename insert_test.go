package sq

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertBuilderToSql(t *testing.T) {
	b := Insert("").
		Prefix("WITH prefix AS ?", 0).
		Into("a").
		Options("DELAYED", "IGNORE").
		Columns("b", "c").
		Values(1, 2).
		Values(3, Expr("? + 1", 4)).
		Suffix("RETURNING ?", 5)

	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSQL :=
		"WITH prefix AS ? " +
			"INSERT DELAYED IGNORE INTO a (b,c) VALUES (?,?),(?,? + 1) " +
			"RETURNING ?"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []any{0, 1, 2, 3, 4, 5}
	require.Equal(t, expectedArgs, args)
}

func TestInsertBuilderToSqlErr(t *testing.T) {
	_, _, err := Insert("").Values(1).ToSql()
	require.Error(t, err)

	_, _, err = Insert("x").ToSql()
	require.Error(t, err)
}

func TestInsertBuilderPlaceholders(t *testing.T) {
	b := Insert("test").Values(1, 2)

	sql, _, _ := b.PlaceholderFormat(Question).ToSql()
	require.Equal(t, "INSERT INTO test VALUES (?,?)", sql)

	sql, _, _ = b.PlaceholderFormat(Dollar).ToSql()
	require.Equal(t, "INSERT INTO test VALUES ($1,$2)", sql)
}

func TestInsertBuilderSetMap(t *testing.T) {
	b := Insert("table").SetMap(Eq{"field1": 1, "field2": 2, "field3": 3})

	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSQL := "INSERT INTO table (field1,field2,field3) VALUES (?,?,?)"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []any{1, 2, 3}
	require.Equal(t, expectedArgs, args)
}

func TestInsertBuilderSelect(t *testing.T) {
	sb := Select("field1").From("table1").Where(Eq{"field1": 1})
	ib := Insert("table2").Columns("field1").Select(sb)

	sql, args, err := ib.ToSql()
	require.NoError(t, err)

	expectedSQL := "INSERT INTO table2 (field1) SELECT field1 FROM table1 WHERE field1 = ?"
	require.Equal(t, expectedSQL, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestInsertBuilderReplace(t *testing.T) {
	b := Replace("table").Values(1)

	expectedSQL := "REPLACE INTO table VALUES (?)"

	sql, _, err := b.ToSql()
	require.NoError(t, err)

	require.Equal(t, expectedSQL, sql)
}
