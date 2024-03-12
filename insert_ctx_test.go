package squirrel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertBuilderContextRunners(t *testing.T) {
	db := &DBStub{}
	b := Insert("test").Values(1).RunWith(db)

	expectedSql := "INSERT INTO test VALUES (?)"

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

func TestInsertBuilderContextNoRunner(t *testing.T) {
	b := Insert("test").Values(1)

	_, err := b.ExecContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	_, err = b.QueryContext(ctx)
	require.Equal(t, RunnerNotSet, err)

	err = b.ScanContext(ctx)
	require.Equal(t, RunnerNotSet, err)
}
