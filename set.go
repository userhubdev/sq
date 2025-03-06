package sq

import (
	"fmt"
)

type setOp struct {
	parts []Sqlizer
	sep   string
}

func (op setOp) ToSql() (string, []any, error) {
	if len(op.parts) == 0 {
		return "", nil, fmt.Errorf("%s has no parts", op.sep)
	}

	b := Builder{}

	for i, p := range op.parts {
		if i > 0 {
			b.WriteString(" ")
			b.WriteString(op.sep)
		}

		b.WriteSql(p)
	}

	return b.ToSql()
}

func UnionAll(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "UNION ALL",
	}
}

func UnionDistinct(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "UNION DISTINCT",
	}
}

func IntersectAll(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "INTERSECT ALL",
	}
}

func IntersectDistinct(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "INTERSECT DISTINCT",
	}
}

func ExceptAll(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "EXCEPT ALL",
	}
}

func ExceptDistinct(parts ...Sqlizer) Sqlizer {
	return setOp{
		parts: parts,
		sep:   "EXCEPT DISTINCT",
	}
}
