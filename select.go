package sq

import (
	"fmt"
	"strings"

	"github.com/userhubdev/sq/internal/builder"
)

type selectData struct {
	PlaceholderFormat PlaceholderFormat
	Prefixes          []Sqlizer
	Options           []string
	Columns           []Sqlizer
	From              Sqlizer
	Joins             []Sqlizer
	WhereParts        []Sqlizer
	GroupBys          []string
	HavingParts       []Sqlizer
	OrderByParts      []Sqlizer
	Limit             string
	Offset            string
	Suffixes          []Sqlizer
}

func (d *selectData) ToSql() (sqlStr string, args []any, err error) {
	sqlStr, args, err = d.ToSqlRaw()
	if err != nil {
		return
	}

	sqlStr, err = d.PlaceholderFormat.ReplacePlaceholders(sqlStr)
	return
}

func (d *selectData) ToSqlRaw() (sqlStr string, args []any, err error) {
	if len(d.Columns) == 0 {
		err = fmt.Errorf("select statements must have at least one result column")
		return
	}

	sql := &strings.Builder{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSql(d.Prefixes, sql, " ", args)
		if err != nil {
			return
		}

		sql.WriteString(" ")
	}

	sql.WriteString("SELECT ")

	if len(d.Options) > 0 {
		sql.WriteString(strings.Join(d.Options, " "))
		sql.WriteString(" ")
	}

	if len(d.Columns) > 0 {
		args, err = appendToSql(d.Columns, sql, ", ", args)
		if err != nil {
			return
		}
	}

	if d.From != nil {
		sql.WriteString(" FROM ")
		args, err = appendToSql([]Sqlizer{d.From}, sql, "", args)
		if err != nil {
			return
		}
	}

	if len(d.Joins) > 0 {
		sql.WriteString(" ")
		args, err = appendToSql(d.Joins, sql, " ", args)
		if err != nil {
			return
		}
	}

	if len(d.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSql(d.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(d.GroupBys, ", "))
	}

	if len(d.HavingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSql(d.HavingParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(d.OrderByParts) > 0 {
		sql.WriteString(" ORDER BY ")
		args, err = appendToSql(d.OrderByParts, sql, ", ", args)
		if err != nil {
			return
		}
	}

	if len(d.Limit) > 0 {
		sql.WriteString(" LIMIT ")
		sql.WriteString(d.Limit)
	}

	if len(d.Offset) > 0 {
		sql.WriteString(" OFFSET ")
		sql.WriteString(d.Offset)
	}

	if len(d.Suffixes) > 0 {
		sql.WriteString(" ")

		args, err = appendToSql(d.Suffixes, sql, " ", args)
		if err != nil {
			return
		}
	}

	sqlStr = sql.String()
	return
}

// Builder

// SelectBuilder builds SQL SELECT statements.
type SelectBuilder builder.Builder

func init() {
	builder.Register(SelectBuilder{}, selectData{})
}

// Format methods

// PlaceholderFormat sets PlaceholderFormat (e.g. Question or Dollar) for the
// query.
func (b SelectBuilder) PlaceholderFormat(f PlaceholderFormat) SelectBuilder {
	return builder.Set(b, "PlaceholderFormat", f).(SelectBuilder)
}

// SQL methods

// ToSql builds the query into a SQL string and bound args.
func (b SelectBuilder) ToSql() (string, []any, error) {
	data := builder.GetStruct(b).(selectData)
	return data.ToSql()
}

func (b SelectBuilder) ToSqlRaw() (string, []any, error) {
	data := builder.GetStruct(b).(selectData)
	return data.ToSqlRaw()
}

// Prefix adds an expression to the beginning of the query
func (b SelectBuilder) Prefix(sql string, args ...any) SelectBuilder {
	return b.PrefixExpr(Expr(sql, args...))
}

// RemovePrefixes removes prefixes.
func (b SelectBuilder) RemovePrefixes() SelectBuilder {
	return builder.Delete(b, "Prefixes").(SelectBuilder)
}

// PrefixExpr adds an expression to the very beginning of the query
func (b SelectBuilder) PrefixExpr(expr Sqlizer) SelectBuilder {
	return builder.Append(b, "Prefixes", expr).(SelectBuilder)
}

// Distinct adds a DISTINCT clause to the query.
func (b SelectBuilder) Distinct() SelectBuilder {
	return b.Options("DISTINCT")
}

// Options adds select option to the query
func (b SelectBuilder) Options(options ...string) SelectBuilder {
	return builder.Extend(b, "Options", options).(SelectBuilder)
}

// Columns adds result columns to the query.
func (b SelectBuilder) Columns(columns ...string) SelectBuilder {
	parts := make([]any, 0, len(columns))
	for _, str := range columns {
		parts = append(parts, newPart(str))
	}
	return builder.Extend(b, "Columns", parts).(SelectBuilder)
}

// RemoveColumns remove all columns from query.
// Must add a new column with Column or Columns methods, otherwise
// return a error.
func (b SelectBuilder) RemoveColumns() SelectBuilder {
	return builder.Delete(b, "Columns").(SelectBuilder)
}

// Column adds a result column to the query.
// Unlike Columns, Column accepts args which will be bound to placeholders in
// the columns string, for example:
//
//	Column("IF(col IN ("+sq.Placeholders(3)+"), 1, 0) as col", 1, 2, 3)
func (b SelectBuilder) Column(column any, args ...any) SelectBuilder {
	return builder.Append(b, "Columns", newPart(column, args...)).(SelectBuilder)
}

// From sets the FROM clause of the query.
func (b SelectBuilder) From(from string) SelectBuilder {
	return builder.Set(b, "From", newPart(from)).(SelectBuilder)
}

// FromSelect sets a subquery into the FROM clause of the query.
func (b SelectBuilder) FromSelect(from Sqlizer, alias string) SelectBuilder {
	return builder.Set(b, "From", Alias(from, alias)).(SelectBuilder)
}

// JoinClause adds a join clause to the query.
func (b SelectBuilder) JoinClause(pred any, args ...any) SelectBuilder {
	return builder.Append(b, "Joins", newPart(pred, args...)).(SelectBuilder)
}

// RemoveJoins removes JOIN clauses.
func (b SelectBuilder) RemoveJoins() SelectBuilder {
	return builder.Delete(b, "Joins").(SelectBuilder)
}

// Join adds a JOIN clause to the query.
func (b SelectBuilder) Join(join string, rest ...any) SelectBuilder {
	return b.JoinClause("JOIN "+join, rest...)
}

// LeftJoin adds a LEFT JOIN clause to the query.
func (b SelectBuilder) LeftJoin(join string, rest ...any) SelectBuilder {
	return b.JoinClause("LEFT JOIN "+join, rest...)
}

// RightJoin adds a RIGHT JOIN clause to the query.
func (b SelectBuilder) RightJoin(join string, rest ...any) SelectBuilder {
	return b.JoinClause("RIGHT JOIN "+join, rest...)
}

// InnerJoin adds a INNER JOIN clause to the query.
func (b SelectBuilder) InnerJoin(join string, rest ...any) SelectBuilder {
	return b.JoinClause("INNER JOIN "+join, rest...)
}

// CrossJoin adds a CROSS JOIN clause to the query.
func (b SelectBuilder) CrossJoin(join string, rest ...any) SelectBuilder {
	return b.JoinClause("CROSS JOIN "+join, rest...)
}

// FullJoin adds a FULL JOIN clause to the query.
func (b SelectBuilder) FullJoin(join string, rest ...any) SelectBuilder {
	return b.JoinClause("FULL JOIN "+join, rest...)
}

// joinOn adds a join on clause to the query,
func (b SelectBuilder) joinOn(prefix string, join any, on Sqlizer) SelectBuilder {
	return b.JoinClause(ConcatExpr(prefix, " ", join, " ON ", on))
}

// JoinOn adds a JOIN ON clause to the query.
func (b SelectBuilder) JoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("JOIN", join, on)
}

// LeftJoinOn adds a LEFT JOIN ON clause to the query.
func (b SelectBuilder) LeftJoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("LEFT JOIN", join, on)
}

// RightJoinOn adds a RIGHT JOIN ON clause to the query.
func (b SelectBuilder) RightJoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("RIGHT JOIN", join, on)
}

// InnerJoinOn adds a INNER JOIN ON clause to the query.
func (b SelectBuilder) InnerJoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("INNER JOIN", join, on)
}

// CrossJoinOn adds a CROSS JOIN ON clause to the query.
func (b SelectBuilder) CrossJoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("CROSS JOIN", join, on)
}

// FullJoinOn adds a FULL JOIN ON clause to the query.
func (b SelectBuilder) FullJoinOn(join any, on Sqlizer) SelectBuilder {
	return b.joinOn("FULL JOIN", join, on)
}

// JoinSelect adds a JOIN ON clause to the query.
func (b SelectBuilder) JoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.JoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// LeftJoinSelect adds a LEFT JOIN ON clause to the query.
func (b SelectBuilder) LeftJoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.LeftJoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// RightJoinSelect adds a RIGHT JOIN ON clause to the query.
func (b SelectBuilder) RightJoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.RightJoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// InnerJoinSelect adds a INNER JOIN ON clause to the query.
func (b SelectBuilder) InnerJoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.InnerJoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// CrossJoinSelect adds a CROSS JOIN ON clause to the query.
func (b SelectBuilder) CrossJoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.CrossJoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// FullJoinSelect adds a FULL JOIN ON clause to the query.
func (b SelectBuilder) FullJoinSelect(join SelectBuilder, alias string, on Sqlizer) SelectBuilder {
	return b.FullJoinOn(Alias(join.PlaceholderFormat(Question), alias), on)
}

// Where adds an expression to the WHERE clause of the query.
//
// Expressions are ANDed together in the generated SQL.
//
// Where accepts several types for its pred argument:
//
// nil OR "" - ignored.
//
// string - SQL expression.
// If the expression has SQL placeholders then a set of arguments must be passed
// as well, one for each placeholder.
//
// map[string]interface{} OR Eq - map of SQL expressions to values. Each key is
// transformed into an expression like "<key> = ?", with the corresponding value
// bound to the placeholder. If the value is nil, the expression will be "<key>
// IS NULL". If the value is an array or slice, the expression will be "<key> IN
// (?,?,...)", with one placeholder for each item in the value. These expressions
// are ANDed together.
//
// Where will panic if pred isn't any of the above types.
func (b SelectBuilder) Where(pred any, args ...any) SelectBuilder {
	if pred == nil || pred == "" {
		return b
	}
	return builder.Append(b, "WhereParts", newWherePart(pred, args...)).(SelectBuilder)
}

// RemoveWhere removes WHERE clause.
func (b SelectBuilder) RemoveWhere() SelectBuilder {
	return builder.Delete(b, "WhereParts").(SelectBuilder)
}

// GroupBy adds GROUP BY expressions to the query.
func (b SelectBuilder) GroupBy(groupBys ...string) SelectBuilder {
	return builder.Extend(b, "GroupBys", groupBys).(SelectBuilder)
}

// RemoveGroupBy removes GROUP BY clause.
func (b SelectBuilder) RemoveGroupBy() SelectBuilder {
	return builder.Delete(b, "GroupBys").(SelectBuilder)
}

// Having adds an expression to the HAVING clause of the query.
//
// See Where.
func (b SelectBuilder) Having(pred any, rest ...any) SelectBuilder {
	return builder.Append(b, "HavingParts", newWherePart(pred, rest...)).(SelectBuilder)
}

// RemoveHaving removes HAVING clause.
func (b SelectBuilder) RemoveHaving() SelectBuilder {
	return builder.Delete(b, "HavingParts").(SelectBuilder)
}

// OrderByClause adds ORDER BY clause to the query.
func (b SelectBuilder) OrderByClause(pred any, args ...any) SelectBuilder {
	return builder.Append(b, "OrderByParts", newPart(pred, args...)).(SelectBuilder)
}

// OrderBy adds ORDER BY expressions to the query.
func (b SelectBuilder) OrderBy(orderBys ...string) SelectBuilder {
	for _, orderBy := range orderBys {
		b = b.OrderByClause(orderBy)
	}

	return b
}

// RemoveOrderBy removes ORDER BY clause.
func (b SelectBuilder) RemoveOrderBy() SelectBuilder {
	return builder.Delete(b, "OrderByParts").(SelectBuilder)
}

// Limit sets a LIMIT clause on the query.
func (b SelectBuilder) Limit(limit uint64) SelectBuilder {
	return builder.Set(b, "Limit", fmt.Sprintf("%d", limit)).(SelectBuilder)
}

// RemoveLimit removes LIMIT clause.
func (b SelectBuilder) RemoveLimit() SelectBuilder {
	return builder.Delete(b, "Limit").(SelectBuilder)
}

// Offset sets a OFFSET clause on the query.
func (b SelectBuilder) Offset(offset uint64) SelectBuilder {
	return builder.Set(b, "Offset", fmt.Sprintf("%d", offset)).(SelectBuilder)
}

// RemoveOffset removes OFFSET clause.
func (b SelectBuilder) RemoveOffset() SelectBuilder {
	return builder.Delete(b, "Offset").(SelectBuilder)
}

// Suffix adds an expression to the end of the query
func (b SelectBuilder) Suffix(sql string, args ...any) SelectBuilder {
	return b.SuffixExpr(Expr(sql, args...))
}

// RemoveSuffixes removes suffixes.
func (b SelectBuilder) RemoveSuffixes() SelectBuilder {
	return builder.Delete(b, "Suffixes").(SelectBuilder)
}

// SuffixExpr adds an expression to the end of the query
func (b SelectBuilder) SuffixExpr(expr Sqlizer) SelectBuilder {
	return builder.Append(b, "Suffixes", expr).(SelectBuilder)
}
