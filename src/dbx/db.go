package dbx

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	DBName string

	db *sql.DB
	//go:embed "sql/*.sql"
	migrations embed.FS
)

func Init() {
	var err error
	if db, err = sql.Open("sqlite3", DBName); err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := applyMigrations("sql"); err != nil {
		log.Fatal(err)
	}

	/* TODO: Remove debug code */
	if err := CreateUser(User{
		Email:    "mail@thomasvoss.com",
		Username: "Thomas",
		Password: "69",
		AdminP:   true,
	}); err != nil {
		log.Fatal(err)
	}
	if err := CreateUser(User{
		Email:    "foo@BAR.baz",
		Username: "Foobar",
		Password: "420",
		AdminP:   false,
	}); err != nil {
		log.Fatal(err)
	}
	if _, err := GetMintages("ad"); err != nil {
		log.Fatal(err)
	}
}

func Close() {
	db.Close()
}

func applyMigrations(dir string) error {
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

	files, err := fs.ReadDir(migrations, dir)
	if err != nil {
		return err
	}

	scripts := []string{}
	for _, f := range files {
		scripts = append(scripts, f.Name())
	}

	sort.Strings(scripts)
	for _, f := range scripts[latest+1:] {
		qry, err := migrations.ReadFile(filepath.Join(dir, f))
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
		log.Printf("Applied database migration ‘%s’", f)
	}

	return nil
}

func scanToStructs[T any](rs *sql.Rows) ([]T, error) {
	xs := []T{}
	for rs.Next() {
		x, err := scanToStruct[T](rs)
		if err != nil {
			return nil, err
		}
		xs = append(xs, x)
	}
	return xs, rs.Err()
}

func scanToStruct[T any](rs *sql.Rows) (T, error) {
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

	rs.Next()
	if err := rs.Err(); err != nil {
		return zero, err
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
