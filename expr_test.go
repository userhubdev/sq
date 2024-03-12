package squirrel

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcatExpr(t *testing.T) {
	b := ConcatExpr("COALESCE(name,", Expr("CONCAT(?,' ',?)", "f", "l"), ")")
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "COALESCE(name,CONCAT(?,' ',?))"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"f", "l"}
	require.Equal(t, expectedArgs, args)
}

func TestConcatExprBadType(t *testing.T) {
	b := ConcatExpr("prefix", 123, "suffix")
	_, _, err := b.ToSql()
	require.Error(t, err)
	require.Contains(t, err.Error(), "123 is not")
}

func TestEqToSql(t *testing.T) {
	b := Eq{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id = ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestEqEmptyToSql(t *testing.T) {
	sql, args, err := Eq{}.ToSql()
	require.NoError(t, err)

	expectedSql := "(1=1)"
	require.Equal(t, expectedSql, sql)
	require.Empty(t, args)
}

func TestEqInToSql(t *testing.T) {
	b := Eq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id IN (?,?,?)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1, 2, 3}
	require.Equal(t, expectedArgs, args)
}

func TestEqSqlizeToSql(t *testing.T) {
	b := Eq{"id": Expr("test + ?", 1)}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id = test + ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)

	b = Eq{"id": Expr("test")}
	sql, args, err = b.ToSql()
	require.NoError(t, err)

	expectedSql = "id = test"
	require.Equal(t, expectedSql, sql)

	require.Empty(t, args)
}

func TestNotEqToSql(t *testing.T) {
	b := NotEq{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id <> ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestEqNotInToSql(t *testing.T) {
	b := NotEq{"id": []int{1, 2, 3}}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id NOT IN (?,?,?)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1, 2, 3}
	require.Equal(t, expectedArgs, args)
}

func TestEqInEmptyToSql(t *testing.T) {
	b := Eq{"id": []int{}}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "(1=0)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{}
	require.Equal(t, expectedArgs, args)
}

func TestNotEqInEmptyToSql(t *testing.T) {
	b := NotEq{"id": []int{}}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "(1=1)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{}
	require.Equal(t, expectedArgs, args)
}

func TestEqBytesToSql(t *testing.T) {
	b := Eq{"id": []byte("test")}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id = ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{[]byte("test")}
	require.Equal(t, expectedArgs, args)
}

func TestLtToSql(t *testing.T) {
	b := Lt{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id < ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestLtSqlizeToSql(t *testing.T) {
	b := Lt{"id": Expr("test + ?", 1)}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id < test + ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestLtOrEqToSql(t *testing.T) {
	b := LtOrEq{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id <= ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestGtToSql(t *testing.T) {
	b := Gt{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id > ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestGtOrEqToSql(t *testing.T) {
	b := GtOrEq{"id": 1}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "id >= ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1}
	require.Equal(t, expectedArgs, args)
}

func TestExprNilToSql(t *testing.T) {
	var b Sqlizer
	b = NotEq{"name": nil}
	sql, args, err := b.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)

	expectedSql := "name IS NOT NULL"
	require.Equal(t, expectedSql, sql)

	b = Eq{"name": nil}
	sql, args, err = b.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)

	expectedSql = "name IS NULL"
	require.Equal(t, expectedSql, sql)
}

func TestNullTypeString(t *testing.T) {
	var b Sqlizer
	var name sql.NullString

	b = Eq{"name": name}
	sql, args, err := b.ToSql()

	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "name IS NULL", sql)

	err = name.Scan("Name")
	require.NoError(t, err)
	b = Eq{"name": name}
	sql, args, err = b.ToSql()

	require.NoError(t, err)
	require.Equal(t, []any{"Name"}, args)
	require.Equal(t, "name = ?", sql)
}

func TestNullTypeInt64(t *testing.T) {
	var userID sql.NullInt64
	err := userID.Scan(nil)
	require.NoError(t, err)
	b := Eq{"user_id": userID}
	sql, args, err := b.ToSql()

	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "user_id IS NULL", sql)

	err = userID.Scan(int64(10))
	require.NoError(t, err)
	b = Eq{"user_id": userID}
	sql, args, err = b.ToSql()

	require.NoError(t, err)
	require.Equal(t, []any{int64(10)}, args)
	require.Equal(t, "user_id = ?", sql)
}

func TestNilPointer(t *testing.T) {
	var name *string = nil
	eq := Eq{"name": name}
	sql, args, err := eq.ToSql()

	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "name IS NULL", sql)

	neq := NotEq{"name": name}
	sql, args, err = neq.ToSql()

	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "name IS NOT NULL", sql)

	var ids *[]int = nil
	eq = Eq{"id": ids}
	sql, args, err = eq.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "id IS NULL", sql)

	neq = NotEq{"id": ids}
	sql, args, err = neq.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "id IS NOT NULL", sql)

	var ida *[3]int = nil
	eq = Eq{"id": ida}
	sql, args, err = eq.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "id IS NULL", sql)

	neq = NotEq{"id": ida}
	sql, args, err = neq.ToSql()
	require.NoError(t, err)
	require.Empty(t, args)
	require.Equal(t, "id IS NOT NULL", sql)

}

func TestNotNilPointer(t *testing.T) {
	c := "Name"
	name := &c
	eq := Eq{"name": name}
	sql, args, err := eq.ToSql()

	require.NoError(t, err)
	require.Equal(t, []any{"Name"}, args)
	require.Equal(t, "name = ?", sql)

	neq := NotEq{"name": name}
	sql, args, err = neq.ToSql()

	require.NoError(t, err)
	require.Equal(t, []any{"Name"}, args)
	require.Equal(t, "name <> ?", sql)

	s := []int{1, 2, 3}
	ids := &s
	eq = Eq{"id": ids}
	sql, args, err = eq.ToSql()
	require.NoError(t, err)
	require.Equal(t, []any{1, 2, 3}, args)
	require.Equal(t, "id IN (?,?,?)", sql)

	neq = NotEq{"id": ids}
	sql, args, err = neq.ToSql()
	require.NoError(t, err)
	require.Equal(t, []any{1, 2, 3}, args)
	require.Equal(t, "id NOT IN (?,?,?)", sql)

	a := [3]int{1, 2, 3}
	ida := &a
	eq = Eq{"id": ida}
	sql, args, err = eq.ToSql()
	require.NoError(t, err)
	require.Equal(t, []any{1, 2, 3}, args)
	require.Equal(t, "id IN (?,?,?)", sql)

	neq = NotEq{"id": ida}
	sql, args, err = neq.ToSql()
	require.NoError(t, err)
	require.Equal(t, []any{1, 2, 3}, args)
	require.Equal(t, "id NOT IN (?,?,?)", sql)
}

func TestEmptyAndToSql(t *testing.T) {
	sql, args, err := And{}.ToSql()
	require.NoError(t, err)

	expectedSql := "(1=1)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{}
	require.Equal(t, expectedArgs, args)
}

func TestEmptyOrToSql(t *testing.T) {
	sql, args, err := Or{}.ToSql()
	require.NoError(t, err)

	expectedSql := "(1=0)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{}
	require.Equal(t, expectedArgs, args)
}

func TestLikeToSql(t *testing.T) {
	b := Like{"name": "%irrel"}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "name LIKE ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"%irrel"}
	require.Equal(t, expectedArgs, args)
}

func TestLikeSqlizeToSql(t *testing.T) {
	b := Like{"name": Expr("CONCAT(test, ?)", "%irrel")}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "name LIKE CONCAT(test, ?)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"%irrel"}
	require.Equal(t, expectedArgs, args)
}

func TestNotLikeToSql(t *testing.T) {
	b := NotLike{"name": "%irrel"}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "name NOT LIKE ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"%irrel"}
	require.Equal(t, expectedArgs, args)
}

func TestILikeToSql(t *testing.T) {
	b := ILike{"name": "sq%"}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "name ILIKE ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"sq%"}
	require.Equal(t, expectedArgs, args)
}

func TestNotILikeToSql(t *testing.T) {
	b := NotILike{"name": "sq%"}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "name NOT ILIKE ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{"sq%"}
	require.Equal(t, expectedArgs, args)
}

func TestSqlEqOrder(t *testing.T) {
	b := Eq{"a": 1, "b": 2, "c": 3}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "a = ? AND b = ? AND c = ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1, 2, 3}
	require.Equal(t, expectedArgs, args)
}

func TestSqlLtOrder(t *testing.T) {
	b := Lt{"a": 1, "b": 2, "c": 3}
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "a < ? AND b < ? AND c < ?"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{1, 2, 3}
	require.Equal(t, expectedArgs, args)
}

func TestExprEscaped(t *testing.T) {
	b := Expr("count(??)", Expr("x"))
	sql, args, err := b.ToSql()
	require.NoError(t, err)

	expectedSql := "count(??)"
	require.Equal(t, expectedSql, sql)

	expectedArgs := []any{Expr("x")}
	require.Equal(t, expectedArgs, args)
}

func TestExprRecursion(t *testing.T) {
	{
		b := Expr("count(?)", Expr("nullif(a,?)", "b"))
		sql, args, err := b.ToSql()
		require.NoError(t, err)

		expectedSql := "count(nullif(a,?))"
		require.Equal(t, expectedSql, sql)

		expectedArgs := []any{"b"}
		require.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("extract(? from ?)", Expr("epoch"), "2001-02-03")
		sql, args, err := b.ToSql()
		require.NoError(t, err)

		expectedSql := "extract(epoch from ?)"
		require.Equal(t, expectedSql, sql)

		expectedArgs := []any{"2001-02-03"}
		require.Equal(t, expectedArgs, args)
	}
	{
		b := Expr("JOIN t1 ON ?", And{Eq{"id": 1}, Expr("NOT c1"), Expr("? @@ ?", "x", "y")})
		sql, args, err := b.ToSql()
		require.NoError(t, err)

		expectedSql := "JOIN t1 ON (id = ? AND NOT c1 AND ? @@ ?)"
		require.Equal(t, expectedSql, sql)

		expectedArgs := []any{1, "x", "y"}
		require.Equal(t, expectedArgs, args)
	}
}

func ExampleEq() {
	Select("id", "created", "first_name").From("users").Where(Eq{
		"company": 20,
	})
}
