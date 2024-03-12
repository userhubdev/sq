package squirrel

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func (s *DBStub) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	s.LastPrepareSql = query
	s.PrepareCount++
	return nil, nil
}

func (s *DBStub) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	s.LastExecSql = query
	s.LastExecArgs = args
	return nil, nil
}

func (s *DBStub) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	s.LastQuerySql = query
	s.LastQueryArgs = args
	return nil, nil
}

func (s *DBStub) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	s.LastQueryRowSql = query
	s.LastQueryRowArgs = args
	return &Row{RowScanner: &RowStub{}}
}

var ctx = context.Background()

func TestExecContextWith(t *testing.T) {
	db := &DBStub{}
	_, err := ExecContextWith(ctx, db, sqlizer)
	require.NoError(t, err)
	require.Equal(t, sqlStr, db.LastExecSql)
}

func TestQueryContextWith(t *testing.T) {
	db := &DBStub{}
	_, err := QueryContextWith(ctx, db, sqlizer)
	require.NoError(t, err)
	require.Equal(t, sqlStr, db.LastQuerySql)
}

func TestQueryRowContextWith(t *testing.T) {
	db := &DBStub{}
	QueryRowContextWith(ctx, db, sqlizer)
	require.Equal(t, sqlStr, db.LastQueryRowSql)
}
