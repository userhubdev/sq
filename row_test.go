package squirrel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type RowStub struct {
	Scanned bool
}

func (r *RowStub) Scan(_ ...any) error {
	r.Scanned = true
	return nil
}

func TestRowScan(t *testing.T) {
	stub := &RowStub{}
	row := &Row{RowScanner: stub}
	err := row.Scan()
	require.True(t, stub.Scanned, "row was not scanned")
	require.NoError(t, err)
}

func TestRowScanErr(t *testing.T) {
	stub := &RowStub{}
	rowErr := fmt.Errorf("scan err")
	row := &Row{RowScanner: stub, err: rowErr}
	err := row.Scan()
	require.False(t, stub.Scanned, "row was scanned")
	require.Equal(t, rowErr, err)
}
