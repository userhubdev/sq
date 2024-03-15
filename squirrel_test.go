package sq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type DBStub struct {
	LastPrepareSql string
	PrepareCount   int

	LastExecSql  string
	LastExecArgs []any

	LastQuerySql  string
	LastQueryArgs []any

	LastQueryRowSql  string
	LastQueryRowArgs []any
}

var testDebugUpdateSQL = Update("table").SetMap(Eq{"x": 1, "y": "val"})
var expectedDebugUpdateSQL = "UPDATE table SET x = '1', y = 'val'"

func TestDebugSqlizerUpdateColon(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugUpdateSQL, DebugSqlizer(testDebugUpdateSQL))
}

func TestDebugSqlizerUpdateAtp(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugUpdateSQL, DebugSqlizer(testDebugUpdateSQL))
}

func TestDebugSqlizerUpdateDollar(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugUpdateSQL, DebugSqlizer(testDebugUpdateSQL))
}

func TestDebugSqlizerUpdateQuestion(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugUpdateSQL, DebugSqlizer(testDebugUpdateSQL))
}

var testDebugDeleteSQL = Delete("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugDeleteSQL = "DELETE FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugSqlizerDeleteColon(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugDeleteSQL, DebugSqlizer(testDebugDeleteSQL))
}

func TestDebugSqlizerDeleteAtp(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugDeleteSQL, DebugSqlizer(testDebugDeleteSQL))
}

func TestDebugSqlizerDeleteDollar(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugDeleteSQL, DebugSqlizer(testDebugDeleteSQL))
}

func TestDebugSqlizerDeleteQuestion(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugDeleteSQL, DebugSqlizer(testDebugDeleteSQL))
}

var testDebugInsertSQL = Insert("table").Values(1, "test")
var expectedDebugInsertSQL = "INSERT INTO table VALUES ('1','test')"

func TestDebugSqlizerInsertColon(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugInsertSQL, DebugSqlizer(testDebugInsertSQL))
}

func TestDebugSqlizerInsertAtp(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugInsertSQL, DebugSqlizer(testDebugInsertSQL))
}

func TestDebugSqlizerInsertDollar(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugInsertSQL, DebugSqlizer(testDebugInsertSQL))
}

func TestDebugSqlizerInsertQuestion(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugInsertSQL, DebugSqlizer(testDebugInsertSQL))
}

var testDebugSelectSQL = Select("*").From("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugSelectSQL = "SELECT * FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugSqlizerSelectColon(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugSelectSQL, DebugSqlizer(testDebugSelectSQL))
}

func TestDebugSqlizerSelectAtp(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugSelectSQL, DebugSqlizer(testDebugSelectSQL))
}

func TestDebugSqlizerSelectDollar(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugSelectSQL, DebugSqlizer(testDebugSelectSQL))
}

func TestDebugSqlizerSelectQuestion(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugSelectSQL, DebugSqlizer(testDebugSelectSQL))
}

func TestDebugSqlizer(t *testing.T) {
	sqlizer := Expr("x = ? AND y = ? AND z = '??'", 1, "text")
	expectedDebug := "x = '1' AND y = 'text' AND z = '?'"
	require.Equal(t, expectedDebug, DebugSqlizer(sqlizer))
}

func TestDebugSqlizerErrors(t *testing.T) {
	errorMsg := DebugSqlizer(Expr("x = ?", 1, 2)) // Not enough placeholders
	require.True(t, strings.HasPrefix(errorMsg, "[DebugSqlizer error: "))

	errorMsg = DebugSqlizer(Expr("x = ? AND y = ?", 1)) // Too many placeholders
	require.True(t, strings.HasPrefix(errorMsg, "[DebugSqlizer error: "))

	errorMsg = DebugSqlizer(Lt{"x": nil}) // Cannot use nil values with Lt
	require.True(t, strings.HasPrefix(errorMsg, "[ToSql error: "))
}
