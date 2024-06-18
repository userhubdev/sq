package sq

import (
	"github.com/userhubdev/sq/internal/builder"
)

func init() {
	builder.Register(WithBuilder{}, withData{})
}

// withPart is a helper structure to describe the cte parts of a WITH clause.
type withPart struct {
	alias string
	cte   Sqlizer
}

// newWithPart creates a new withPart for a WITH clause.
func newWithPart(alias string, cte Sqlizer) withPart {
	return withPart{alias: alias, cte: cte}
}

// withData holds all the data required to build a WITH clause.
type withData struct {
	PlaceholderFormat PlaceholderFormat
	WithParts         []withPart
}

// ToSql implements Sqlizer.
func (d *withData) ToSql() (sqlStr string, args []any, err error) {
	if len(d.WithParts) == 0 {
		return "", nil, nil
	}

	sql := Builder{}

	sql.WriteString("WITH")

	for i, p := range d.WithParts {
		if i > 0 {
			sql.WriteString(", ")
		} else {
			sql.WriteString(" ")
		}
		sql.WriteString(p.alias)
		sql.WriteString(" AS (")
		sql.WriteSql(p.cte)
		sql.WriteString(")")
	}

	return sql.ToSql()
}

// WithBuilder builds a WITH clause.
type WithBuilder builder.Builder

// ToSql builds the query into a SQL string and bound args.
func (b WithBuilder) ToSql() (string, []any, error) {
	data := builder.GetStruct(b).(withData)
	return data.ToSql()
}

// As adds a "... AS (...)" part to the WITH clause.
func (b WithBuilder) As(alias string, sql Sqlizer) WithBuilder {
	return builder.Append(b, "WithParts", newWithPart(alias, sql)).(WithBuilder)
}

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// WITH clause.
func (b WithBuilder) PlaceholderFormat(f PlaceholderFormat) WithBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(WithBuilder)
}

// Select starts a primary SELECT statement for the WITH clause.
func (b WithBuilder) Select(columns ...string) SelectBuilder {
	data := builder.GetStruct(b).(withData)

	sql := StatementBuilder.Select(columns...)

	if len(data.WithParts) > 0 {
		sql = sql.PrefixExpr(&data)
	}

	return sql.PlaceholderFormat(data.PlaceholderFormat)
}

// Insert starts a primary INSERT statement for the WITH clause.
func (b WithBuilder) Insert(into string) InsertBuilder {
	data := builder.GetStruct(b).(withData)

	sql := StatementBuilder.Insert(into)

	if len(data.WithParts) > 0 {
		sql = sql.PrefixExpr(&data)
	}

	return sql.PlaceholderFormat(data.PlaceholderFormat)

}

// Update starts a primary UPDATE statement for the WITH clause.
func (b WithBuilder) Update(table string) UpdateBuilder {
	data := builder.GetStruct(b).(withData)

	sql := StatementBuilder.Update(table)

	if len(data.WithParts) > 0 {
		sql = sql.PrefixExpr(&data)
	}

	return sql.PlaceholderFormat(data.PlaceholderFormat)
}

// Delete starts a primary DELETE statement for the WITH clause.
func (b WithBuilder) Delete(from string) DeleteBuilder {
	data := builder.GetStruct(b).(withData)

	sql := StatementBuilder.Delete(from)

	if len(data.WithParts) > 0 {
		sql = sql.PrefixExpr(&data)
	}

	return sql.PlaceholderFormat(data.PlaceholderFormat)
}
