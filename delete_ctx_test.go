package squirrel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteBuilderContextRunners(t *testing.T) {
	db := &DBStub{}
	b := Delete("test").Where("x = ?", 1).RunWith(db)

	expectedSql := "DELETE FROM test WHERE x = ?"

	_, err := b.ExecContext(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedSql, db.LastExecSql)

	_, err = b.QueryContext(ctx)
	require.NoError(t, err)
	require.Equal(t, expectedSql, db.LastQuerySql)

	b.QueryRowContext(ctx)
	require.Equal(t, expectedSql, db.LastQueryRowSql)

	err = b.ScanContext(ctx)
	require.NoError(t, err)
}

func TestDeleteBuilderContextNoRunner(t *testing.T) {
	b := Delete("test").Where("x != ?", 0).Suffix("RETURNING x")

	_, err := b.ExecContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	_, err = b.QueryContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	err = b.ScanContext(ctx)
	require.Equal(t, RunnerNotSet, err)
}
