module git.thomasvoss.com/euro-cash.eu

go 1.24

require (
	github.com/mattn/go-sqlite3 v1.14.28
	golang.org/x/crypto v0.39.0
	golang.org/x/text v0.26.0
	golang.org/x/tools v0.34.0
)

require (
	golang.org/x/mod v0.25.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
)

tool (
	git.thomasvoss.com/euro-cash.eu
	golang.org/x/text/cmd/gotext
)
