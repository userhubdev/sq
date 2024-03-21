package sq

import "strings"

// Builder is a helper that allows to write many Sqlizers one by one
// without constant checks for errors that may come from Sqlizer
type Builder struct {
	strings.Builder
	args []any
	err  error
}

// WriteSql converts Sqlizer to SQL strings and writes it to strings.Builder
func (b *Builder) WriteSql(item Sqlizer) {
	if b.err != nil {
		return
	}

	var str string
	var args []any
	str, args, b.err = nestedToSql(item)

	if b.err != nil {
		return
	}

	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(str)
	b.args = append(b.args, args...)
}

func (b *Builder) ToSql() (string, []any, error) {
	if b.err != nil {
		return "", nil, b.err
	}
	return b.String(), b.args, nil
}
