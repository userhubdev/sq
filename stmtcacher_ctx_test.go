package squirrel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStmtCacherPrepareContext(t *testing.T) {
	db := &DBStub{}
	sc := NewStmtCache(db)
	query := "SELECT 1"

	_, err := sc.PrepareContext(ctx, query)
	require.NoError(t, err)
	require.Equal(t, query, db.LastPrepareSql)

	_, err = sc.PrepareContext(ctx, query)
	require.NoError(t, err)
	require.Equal(t, 1, db.PrepareCount, "expected 1 Prepare, got %d", db.PrepareCount)
}
