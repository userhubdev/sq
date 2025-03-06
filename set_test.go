package sq

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetOp(t *testing.T) {
	testCases := []struct {
		Sep string
		Fn  func(...Sqlizer) Sqlizer
	}{
		{"UNION ALL", UnionAll},
		{"UNION DISTINCT", UnionDistinct},
		{"INTERSECT ALL", IntersectAll},
		{"INTERSECT DISTINCT", IntersectDistinct},
		{"EXCEPT ALL", ExceptAll},
		{"EXCEPT DISTINCT", ExceptDistinct},
	}

	for _, tc := range testCases {
		t.Run(tc.Sep, func(t *testing.T) {
			_, _, err := tc.Fn().ToSql()
			require.EqualError(t, err, fmt.Sprintf("%s has no parts", tc.Sep))

			sql, args, err := tc.Fn(Select("1").Column("?", 2)).ToSql()
			require.NoError(t, err)
			require.Equal(t, "SELECT 1, ?", sql)
			require.Equal(t, []any{2}, args)

			sql, args, err = tc.Fn(
				Select().Column("?", 1),
				Select().Column("?", 2),
			).ToSql()
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("SELECT ? %s SELECT ?", tc.Sep), sql)
			require.Equal(t, []any{1, 2}, args)

			sql, args, err = tc.Fn(
				Select("1").From("a"),
				Select("2").From("b"),
				Select("3").From("c"),
			).ToSql()
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("SELECT 1 FROM a %s SELECT 2 FROM b %s SELECT 3 FROM c", tc.Sep, tc.Sep), sql)
			require.Empty(t, args)

			sql, args, err = Select("*").FromSelect(tc.Fn(
				Select().Column("?", 1),
				Select().Column("?", 2),
			), "s").ToSql()
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("SELECT * FROM (SELECT ? %s SELECT ?) AS s", tc.Sep), sql)
			require.Equal(t, []any{1, 2}, args)
		})
	}
}
