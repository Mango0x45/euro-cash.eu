package dbx

import (
	"fmt"
	"io/fs"
	"log"
	"sort"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

var (
	db     *sqlx.DB
	DBName string
)

func Init(sqlDir fs.FS) {
	db = sqlx.MustConnect("sqlite3", DBName)
	atexit.Register(Close)
	Try(applyMigrations(sqlDir))

	/* TODO: Remove debug code */
	/* Try(CreateUser(User{
	   	Email:    "mail@thomasvoss.com",
	   	Username: "Thomas",
	   	Password: "69",
	   	AdminP:   true,
	   }))
	   Try(CreateUser(User{
	   	Email:    "foo@BAR.baz",
	   	Username: "Foobar",
	   	Password: "420",
	   	AdminP:   false,
	   }))
	   Try2(GetMintages("ad", TypeCirc)) */
}

func Close() {
	db.Close()
}

func applyMigrations(dir fs.FS) error {
	var latest int
	migratedp := true

	err := db.QueryRow("SELECT latest FROM migration").Scan(&latest)
	if err != nil {
		e, ok := err.(sqlite3.Error)
		/* IDK if there is a better way to do this… lol */
		if ok && e.Error() == "no such table: migration" {
			migratedp = false
		} else {
			return err
		}
	}

	if !migratedp {
		latest = -1
	}

	files, err := fs.ReadDir(dir, ".")
	if err != nil {
		return err
	}

	var (
		last    string
		scripts []string
	)

	for _, f := range files {
		if n := f.Name(); n == "last.sql" {
			last = n
		} else {
			scripts = append(scripts, f.Name())
		}
	}

	sort.Strings(scripts)
	for _, f := range scripts[latest+1:] {
		qry, err := fs.ReadFile(dir, f)
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		var n int
		if _, err = fmt.Sscanf(f, "%d", &n); err != nil {
			goto error
		}

		if _, err = tx.Exec(string(qry)); err != nil {
			err = fmt.Errorf("error in ‘%s’: %w", f, err)
			goto error
		}

		_, err = tx.Exec("UPDATE migration SET latest = ? WHERE id = 1", n)
		if err != nil {
			goto error
		}

		if err = tx.Commit(); err != nil {
			goto error
		}

		log.Printf("Applied database migration ‘%s’\n", f)
		continue

	error:
		tx.Rollback()
		return err
	}

	if last != "" {
		qry, err := fs.ReadFile(dir, last)
		if err != nil {
			return err
		}
		if _, err := db.Exec(string(qry)); err != nil {
			return fmt.Errorf("error in ‘%s’: %w", last, err)
		}
		log.Printf("Ran ‘%s’\n", last)
	}

	return nil
}
