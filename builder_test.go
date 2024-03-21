package sq

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	b := Builder{}
	b.WriteSql(Raw("test"))
	sql, args, err := b.ToSql()
	require.Equal(t, "test", sql)
	require.Empty(t, args)
	require.NoError(t, err)

	b = Builder{}
	b.WriteSql(Raw("one"))
	b.WriteSql(Raw("two"))
	sql, args, err = b.ToSql()
	require.Equal(t, "one two", sql)
	require.Empty(t, args)
	require.NoError(t, err)

	b = Builder{}
	b.WriteString("one")
	b.WriteSql(Raw("two"))
	sql, args, err = b.ToSql()
	require.Equal(t, "one two", sql)
	require.Empty(t, args)
	require.NoError(t, err)

	b = Builder{}
	b.WriteSql(Expr("one = ?", 1))
	sql, args, err = b.ToSql()
	require.Equal(t, "one = ?", sql)
	require.Len(t, args, 1)
	require.Equal(t, 1, args[0])
	require.NoError(t, err)

	b = Builder{}
	b.WriteSql(Expr("one = ?", 1))
	b.WriteSql(Expr("AND two = ?", 2))
	sql, args, err = b.ToSql()
	require.Equal(t, "one = ? AND two = ?", sql)
	require.Len(t, args, 2)
	require.Equal(t, 1, args[0])
	require.Equal(t, 2, args[1])
	require.NoError(t, err)

	b = Builder{}
	b.WriteSql(Expr("one = ?", 1))
	b.WriteSql(sqlizeError{err: errors.New("fail")})
	b.WriteSql(Expr("two = ?", 2))
	sql, args, err = b.ToSql()
	require.Empty(t, sql)
	require.Empty(t, args)
	require.EqualError(t, err, "fail")
}

type sqlizeError struct {
	err error
}

func (se sqlizeError) ToSql() (string, []any, error) { return "", nil, se.err }
