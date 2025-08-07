package dbx

import (
	"context"
	"database/sql"
	"slices"
)

type CountryMintageData struct {
	Standard      []MSCountryRow
	Commemorative []MCommemorative
}

type YearMintageData struct {
	Standard      []MSYearRow
	Commemorative []MCommemorative
}

type msRow struct {
	Country      string
	Type         MintageType
	Year         int
	Denomination float64
	Mintmark     sql.Null[string]
	Mintage      sql.Null[int]
	Reference    sql.Null[string]
}

type MSCountryRow struct {
	Year       int
	Mintmark   sql.Null[string]
	Mintages   [ndenoms]sql.Null[int]
	References [ndenoms]sql.Null[string]
}

type MSYearRow struct {
	Country    string
	Mintmark   sql.Null[string]
	Mintages   [ndenoms]sql.Null[int]
	References [ndenoms]sql.Null[string]
}

type MCommemorative struct {
	Country   string
	Type      MintageType
	Year      int
	Name      string
	Number    int
	Mintmark  sql.Null[string]
	Mintage   sql.Null[int]
	Reference sql.Null[string]
}

type MintageType int

/* DO NOT REORDER! */
const (
	TypeCirc MintageType = iota
	TypeNifc
	TypeProof
)

const ndenoms = 8

func NewMintageType(s string) MintageType {
	switch s {
	case "circ":
		return TypeCirc
	case "nifc":
		return TypeNifc
	case "proof":
		return TypeProof
	}
	/* We can get here if the user sends a request manually, so just
	   fallback to this */
	return TypeCirc
}

func GetMintagesByYear(year int, typ MintageType) (YearMintageData, error) {
	var (
		zero YearMintageData
		xs   []MSYearRow
		ys   []MCommemorative
	)

	rs, err := db.QueryxContext(context.TODO(), `
		SELECT * FROM mintages_s
		WHERE year = ? AND type = ?
		ORDER BY country, mintmark, denomination
	`, year, typ)
	if err != nil {
		return zero, err
	}

	for rs.Next() {
		var x msRow
		if err = rs.StructScan(&x); err != nil {
			return zero, err
		}

	loop:
		msr := MSYearRow{
			Country:  x.Country,
			Mintmark: x.Mintmark,
		}
		i := denomToIdx(x.Denomination)
		msr.Mintages[i] = x.Mintage
		msr.References[i] = x.Reference

		for rs.Next() {
			var y msRow
			if err = rs.StructScan(&y); err != nil {
				return zero, err
			}

			if x.Country != y.Country || x.Mintmark != y.Mintmark {
				x = y
				xs = append(xs, msr)
				goto loop
			}

			i = denomToIdx(y.Denomination)
			msr.Mintages[i] = y.Mintage
			msr.References[i] = y.Reference
		}

		xs = append(xs, msr)
	}

	if err = rs.Err(); err != nil {
		return zero, err
	}

	db.SelectContext(context.TODO(), &ys, `
		SELECT * FROM mintages_c
		WHERE year = ? and type = ?
		ORDER BY country, mintmark, number
	`, year, typ)

	return YearMintageData{xs, ys}, nil
}

func GetMintagesByCountry(country string, typ MintageType) (CountryMintageData, error) {
	var (
		zero CountryMintageData
		xs   []MSCountryRow
		ys   []MCommemorative
	)

	rs, err := db.QueryxContext(context.TODO(), `
		SELECT * FROM mintages_s
 		WHERE country = ? AND type = ?
		ORDER BY year, mintmark, denomination
	`, country, typ)
	if err != nil {
		return zero, err
	}

	for rs.Next() {
		var x msRow
		if err = rs.StructScan(&x); err != nil {
			return zero, err
		}

	loop:
		msr := MSCountryRow{
			Year:     x.Year,
			Mintmark: x.Mintmark,
		}
		i := denomToIdx(x.Denomination)
		msr.Mintages[i] = x.Mintage
		msr.References[i] = x.Reference

		for rs.Next() {
			var y msRow
			if err = rs.StructScan(&y); err != nil {
				return zero, err
			}

			if x.Year != y.Year || x.Mintmark != y.Mintmark {
				x = y
				xs = append(xs, msr)
				goto loop
			}

			i = denomToIdx(y.Denomination)
			msr.Mintages[i] = y.Mintage
			msr.References[i] = y.Reference
		}

		xs = append(xs, msr)
	}

	if err = rs.Err(); err != nil {
		return zero, err
	}

	db.SelectContext(context.TODO(), &ys, `
		SELECT * FROM mintages_c
		WHERE country = ? and type = ?
		ORDER BY year, mintmark, number
	`, country, typ)

	return CountryMintageData{xs, ys}, rs.Err()
}

func denomToIdx(d float64) int {
	return slices.Index([]float64{
		0.01, 0.02, 0.05, 0.10,
		0.20, 0.50, 1.00, 2.00,
	}, d)
}
