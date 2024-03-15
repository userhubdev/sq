> **Note**
> This is a fork of [Masterminds/squirrel](https://github.com/Masterminds/squirrel).

# Squirrel - fluent SQL generator for Go

```go
import "github.com/userhubdev/sq"
```

[![GoDoc](https://pkg.go.dev/badge/github.com/userhubdev/squirrel)](https://pkg.go.dev/github.com/userhubdev/squirrel)

**Squirrel is not an ORM.** For an application of Squirrel, check out
[structable, a table-struct mapper](https://github.com/Masterminds/structable)


Squirrel helps you build SQL queries from composable parts:

```go
import "github.com/userhubdev/sq"

users := sq.Select("*").From("users").Join("emails USING (email_id)")

active := users.Where(sq.Eq{"deleted_at": nil})

sql, args, err := active.ToSql()

sql == "SELECT * FROM users JOIN emails USING (email_id) WHERE deleted_at IS NULL"
```

```go
sql, args, err := sq.
    Insert("users").Columns("name", "age").
    Values("moe", 13).Values("larry", sq.Expr("? + 5", 12)).
    ToSql()

sql == "INSERT INTO users (name,age) VALUES (?,?),(?,? + 5)"
```

Squirrel makes conditional query building a breeze:

```go
if len(q) > 0 {
    users = users.Where("name LIKE ?", fmt.Sprint("%", q, "%"))
}
```

## License

Squirrel is released under the
[MIT License](http://www.opensource.org/licenses/MIT).
