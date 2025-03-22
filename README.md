
# Web Url shortener

This project was made on purpose of learning [Go](https://go.dev/). It's probably really unoptimized, any criticizm is welcome.
It allows users to enter a long Url, which will then be wrapped into a shorter one.

## Prerequisites
- [SQLite3](https://www.sqlite.org/)

## Usage

To run it, from projects directory do

```bash
go run .
```
## Packages used

 - [database/sql - provides a generic interface around SQL (or SQL-like) databases](https://pkg.go.dev/database/sql)
 -	[html/template                    - implements data-driven templates for generating HTML output safe against code injection](https://pkg.go.dev/html/template)
 -	[log                - implements a simple logging package](https://pkg.go.dev/log)
 -	[math/rand         - implements pseudo-random number generators suitable for tasks such as simulation](https://pkg.go.dev/math/rand)
 -	[net/http - provides HTTP client and server implementations](https://pkg.go.dev/net/http)
 -	[time - provides functionality for measuring and displaying time](https://pkg.go.dev/time)
 -	[github.com/mattn/go-sqlite3 - sqlite3 driver that conforms to the built-in database/sql interface](https://github.com/mattn/go-sqlite3)
