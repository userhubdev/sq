package sq

import (
	"errors"

	"github.com/userhubdev/sq/internal/builder"
)

func init() {
	builder.Register(CaseBuilder{}, caseData{})
}

// whenPart is a helper structure to describe SQLs "WHEN ... THEN ..." expression
type whenPart struct {
	when Sqlizer
	then Sqlizer
}

func newWhenPart(when any, then any) whenPart {
	return whenPart{newPart(when), newPart(then)}
}

// caseData holds all the data required to build a CASE SQL construct
type caseData struct {
	What      Sqlizer
	WhenParts []whenPart
	Else      Sqlizer
}

// ToSql implements Sqlizer
func (d *caseData) ToSql() (sqlStr string, args []any, err error) {
	if len(d.WhenParts) == 0 {
		err = errors.New("case expression must contain at lease one WHEN clause")

		return
	}

	sql := Builder{}

	sql.WriteString("CASE")
	if d.What != nil {
		sql.WriteSql(d.What)
	}

	for _, p := range d.WhenParts {
		sql.WriteString(" WHEN")
		sql.WriteSql(p.when)
		sql.WriteString(" THEN")
		sql.WriteSql(p.then)
	}

	if d.Else != nil {
		sql.WriteString(" ELSE")
		sql.WriteSql(d.Else)
	}

	sql.WriteString(" END")

	return sql.ToSql()
}

// CaseBuilder builds SQL CASE construct which could be used as parts of queries.
type CaseBuilder builder.Builder

// ToSql builds the query into a SQL string and bound args.
func (b CaseBuilder) ToSql() (string, []any, error) {
	data := builder.GetStruct(b).(caseData)
	return data.ToSql()
}

// what sets optional value for CASE construct "CASE [value] ..."
func (b CaseBuilder) what(expr any) CaseBuilder {
	return builder.Set(b, "What", newPart(expr)).(CaseBuilder)
}

// When adds "WHEN ... THEN ..." part to CASE construct
func (b CaseBuilder) When(when any, then any) CaseBuilder {
	// TODO: performance hint: replace slice of WhenPart with just slice of parts
	// where even indices of the slice belong to "when"s and odd indices belong to "then"s
	return builder.Append(b, "WhenParts", newWhenPart(when, then)).(CaseBuilder)
}

// What sets optional "ELSE ..." part for CASE construct
func (b CaseBuilder) Else(expr any) CaseBuilder {
	return builder.Set(b, "Else", newPart(expr)).(CaseBuilder)
}
