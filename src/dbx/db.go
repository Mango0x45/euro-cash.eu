package dbx

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"reflect"
	"sort"
	"strings"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
	"github.com/mattn/go-sqlite3"
)

var (
	db     *sql.DB
	DBName string
)

func Init(sqlDir fs.FS) {
	db = Try2(sql.Open("sqlite3", DBName))
	Try(db.Ping())
	atexit.Register(Close)
	Try(applyMigrations(sqlDir))

	/* TODO: Remove debug code */
	Try(CreateUser(User{
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
	Try2(GetMintages("ad"))
}

func Close() {
	db.Close()
}

func applyMigrations(dir fs.FS) error {
	var latest int
	migratedp := true

	rows, err := db.Query("SELECT latest FROM migration")
	if err != nil {
		e, ok := err.(sqlite3.Error)
		/* IDK if there is a better way to do this… lol */
		if ok && e.Error() == "no such table: migration" {
			migratedp = false
		} else {
			return err
		}
	} else {
		defer rows.Close()
	}

	if migratedp {
		rows.Next()
		if err := rows.Err(); err != nil {
			return err
		}
		if err := rows.Scan(&latest); err != nil {
			return err
		}
	} else {
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

		if _, err := tx.Exec(string(qry)); err != nil {
			tx.Rollback()
			return fmt.Errorf("error in ‘%s’: %w", f, err)
		}

		var n int
		if _, err := fmt.Sscanf(f, "%d", &n); err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE migration SET latest = ? WHERE id = 1", n)
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		log.Printf("Applied database migration ‘%s’\n", f)
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

func scanToStruct[T any](rs *sql.Rows) (T, error) {
	return scanToStruct2[T](rs, true)
}

func scanToStructs[T any](rs *sql.Rows) ([]T, error) {
	xs := []T{}
	for rs.Next() {
		x, err := scanToStruct2[T](rs, false)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, rs.Err()
}

func scanToStruct2[T any](rs *sql.Rows, callNextP bool) (T, error) {
	var t, zero T

	cols, err := rs.Columns()
	if err != nil {
		return zero, err
	}

	v := reflect.ValueOf(&t).Elem()
	tType := v.Type()

	rawValues := make([]any, len(cols))
	for i := range rawValues {
		var zero any
		rawValues[i] = &zero
	}

	if callNextP {
		rs.Next()
		if err := rs.Err(); err != nil {
			return zero, err
		}
	}
	if err := rs.Scan(rawValues...); err != nil {
		return zero, err
	}

	/* col idx → [field idx, array idx] */
	arrayTargets := make(map[int][2]int)
	colToField := make(map[string]int)

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}

		if strings.Contains(tag, ";") {
			dbcols := strings.Split(tag, ";")
			fv := v.Field(i)
			if fv.Kind() != reflect.Array {
				return zero, fmt.Errorf("field ‘%s’ is not array",
					field.Name)
			}
			if len(dbcols) != fv.Len() {
				return zero, fmt.Errorf("field ‘%s’ array length mismatch",
					field.Name)
			}
			for j, colName := range cols {
				for k, dbColName := range dbcols {
					if colName == dbColName {
						arrayTargets[j] = [2]int{i, k}
					}
				}
			}
		} else {
			colToField[tag] = i
		}
	}

	for i, col := range cols {
		vp := rawValues[i].(*any)
		if fieldIdx, ok := colToField[col]; ok {
			assignValue(v.Field(fieldIdx), *vp)
		} else if target, ok := arrayTargets[i]; ok {
			assignValue(v.Field(target[0]).Index(target[1]), *vp)
		}
	}

	return t, nil
}

func assignValue(fv reflect.Value, val any) {
	if val == nil {
		fv.Set(reflect.Zero(fv.Type()))
		return
	}
	v := reflect.ValueOf(val)
	if v.Type().ConvertibleTo(fv.Type()) {
		fv.Set(v.Convert(fv.Type()))
	}
}
