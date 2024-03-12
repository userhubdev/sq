package squirrel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStmtCachePrepare(t *testing.T) {
	db := &DBStub{}
	sc := NewStmtCache(db)
	query := "SELECT 1"

	_, err := sc.Prepare(query)
	require.NoError(t, err)
	require.Equal(t, query, db.LastPrepareSql)

	_, err = sc.Prepare(query)
	require.NoError(t, err)
	require.Equal(t, 1, db.PrepareCount, "expected 1 Prepare, got %d", db.PrepareCount)

	// clear statement cache
	require.Nil(t, sc.Clear())

	// should prepare the query again
	_, err = sc.Prepare(query)
	require.NoError(t, err)
	require.Equal(t, 2, db.PrepareCount, "expected 2 Prepare, got %d", db.PrepareCount)
}
