package sq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testDebugUpdateSQL = Update("table").SetMap(Eq{"x": 1, "y": "val"})
var expectedDebugUpdateSQL = "UPDATE table SET x = '1', y = 'val'"

func TestDebugUpdateColon(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugUpdateSQL, Debug(testDebugUpdateSQL))
}

func TestDebugUpdateAtp(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugUpdateSQL, Debug(testDebugUpdateSQL))
}

func TestDebugUpdateDollar(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugUpdateSQL, Debug(testDebugUpdateSQL))
}

func TestDebugUpdateQuestion(t *testing.T) {
	testDebugUpdateSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugUpdateSQL, Debug(testDebugUpdateSQL))
}

var testDebugDeleteSQL = Delete("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugDeleteSQL = "DELETE FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugDeleteColon(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugDeleteSQL, Debug(testDebugDeleteSQL))
}

func TestDebugDeleteAtp(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugDeleteSQL, Debug(testDebugDeleteSQL))
}

func TestDebugDeleteDollar(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugDeleteSQL, Debug(testDebugDeleteSQL))
}

func TestDebugDeleteQuestion(t *testing.T) {
	testDebugDeleteSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugDeleteSQL, Debug(testDebugDeleteSQL))
}

var testDebugInsertSQL = Insert("table").Values(1, "test")
var expectedDebugInsertSQL = "INSERT INTO table VALUES ('1','test')"

func TestDebugInsertColon(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugInsertSQL, Debug(testDebugInsertSQL))
}

func TestDebugInsertAtp(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugInsertSQL, Debug(testDebugInsertSQL))
}

func TestDebugInsertDollar(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugInsertSQL, Debug(testDebugInsertSQL))
}

func TestDebugInsertQuestion(t *testing.T) {
	testDebugInsertSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugInsertSQL, Debug(testDebugInsertSQL))
}

var testDebugSelectSQL = Select("*").From("table").Where(And{
	Eq{"column": "val"},
	Eq{"other": 1},
})
var expectedDebugSelectSQL = "SELECT * FROM table WHERE (column = 'val' AND other = '1')"

func TestDebugSelectColon(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Colon)
	require.Equal(t, expectedDebugSelectSQL, Debug(testDebugSelectSQL))
}

func TestDebugSelectAtp(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(AtP)
	require.Equal(t, expectedDebugSelectSQL, Debug(testDebugSelectSQL))
}

func TestDebugSelectDollar(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Dollar)
	require.Equal(t, expectedDebugSelectSQL, Debug(testDebugSelectSQL))
}

func TestDebugSelectQuestion(t *testing.T) {
	testDebugSelectSQL.PlaceholderFormat(Question)
	require.Equal(t, expectedDebugSelectSQL, Debug(testDebugSelectSQL))
}

func TestDebug(t *testing.T) {
	sqlizer := Expr("x = ? AND y = ? AND z = '??'", 1, "text")
	expectedDebug := "x = '1' AND y = 'text' AND z = '?'"
	require.Equal(t, expectedDebug, Debug(sqlizer))
}

func TestDebugErrors(t *testing.T) {
	errorMsg := Debug(Expr("x = ?", 1, 2)) // Not enough placeholders
	require.True(t, strings.HasPrefix(errorMsg, "[Debug error: "))

	errorMsg = Debug(Expr("x = ? AND y = ?", 1)) // Too many placeholders
	require.True(t, strings.HasPrefix(errorMsg, "[Debug error: "))

	errorMsg = Debug(Lt{"x": nil}) // Cannot use nil values with Lt
	require.True(t, strings.HasPrefix(errorMsg, "[ToSql error: "))
}
