module github.com/userhubdev/squirrel/integration

go 1.22

toolchain go1.22.1

replace github.com/userhubdev/squirrel => ../

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/stretchr/testify v1.9.0
	github.com/userhubdev/squirrel v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
