package squirrel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectBuilderContextRunners(t *testing.T) {
	db := &DBStub{}
	b := Select("test").RunWith(db)

	expectedSql := "SELECT test"

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

func TestSelectBuilderContextNoRunner(t *testing.T) {
	b := Select("test")

	_, err := b.ExecContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	_, err = b.QueryContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	err = b.ScanContext(ctx)
	require.Equal(t, RunnerNotSet, err)
}
