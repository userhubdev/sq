module github.com/userhubdev/squirrel/integration

go 1.19

replace github.com/userhubdev/squirrel => ../

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/stretchr/testify v1.8.1
	github.com/userhubdev/squirrel v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
