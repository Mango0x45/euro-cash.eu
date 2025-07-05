package dbx

import (
	"database/sql"
	"slices"
)

type MintageData struct {
	Standard      []MSRow
	Commemorative []MCRow
}

type msRowInternal struct {
	Country      string
	Type         MintageType
	Year         int
	Denomination float64
	Mintmark     sql.Null[string]
	Mintage      sql.Null[int]
	Reference    sql.Null[string]
}

type mcRowInternal struct {
	Country   string
	Type      MintageType
	Year      int
	Name      string
	Number    int
	Mintmark  sql.Null[string]
	Mintage   sql.Null[int]
	Reference sql.Null[string]
}

type MSRow struct {
	Year       int
	Mintmark   string
	Mintages   [ndenoms]int
	References []string
}

type MCRow struct {
	Year      int
	Name      string
	Number    int
	Mintmark  string
	Mintage   int
	Reference string
}

type MintageType int

/* DO NOT REORDER! */
const (
	TypeCirc MintageType = iota
	TypeNifc
	TypeProof
)

/* DO NOT REORDER! */
const (
	MintageUnknown = -iota - 1
	MintageInvalid
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
	/* TODO: Handle this */
	panic("TODO")
}

func GetMintages(country string, typ MintageType) (MintageData, error) {
	var (
		zero MintageData
		xs   []MSRow
		ys   []MCRow
	)

	rs, err := db.Queryx(`
		SELECT * FROM mintages_s
 		WHERE country = ? AND type = ?
		ORDER BY year, mintmark, denomination
	`, country, typ)
	if err != nil {
		return zero, err
	}

	for rs.Next() {
		var x msRowInternal
		if err = rs.StructScan(&x); err != nil {
			return zero, err
		}

	loop:
		msr := MSRow{
			Year:       x.Year,
			Mintmark:   sqlOr(x.Mintmark, ""),
			References: make([]string, 0, ndenoms),
		}
		for i := range msr.Mintages {
			msr.Mintages[i] = MintageUnknown
		}
		msr.Mintages[denomToIdx(x.Denomination)] =
			sqlOr(x.Mintage, MintageUnknown)
		if x.Reference.Valid {
			msr.References = append(msr.References, x.Reference.V)
		}

		for rs.Next() {
			var y msRowInternal
			if err = rs.StructScan(&y); err != nil {
				return zero, err
			}

			if x.Year != y.Year || x.Mintmark != y.Mintmark {
				x = y
				xs = append(xs, msr)
				goto loop
			}

			msr.Mintages[denomToIdx(y.Denomination)] =
				sqlOr(y.Mintage, MintageUnknown)
			if y.Reference.Valid {
				msr.References = append(msr.References, y.Reference.V)
			}
		}

		xs = append(xs, msr)
	}

	if err = rs.Err(); err != nil {
		return zero, err
	}

	rs, err = db.Queryx(`
	   	SELECT * FROM mintages_c
 	   	WHERE country = ? AND type = ?
	   	ORDER BY year, mintmark, number
	`, country, typ)
	if err != nil {
		return zero, err
	}

	for rs.Next() {
		var y mcRowInternal
		if err = rs.StructScan(&y); err != nil {
			return zero, err
		}
		ys = append(ys, MCRow{
			Year:      y.Year,
			Name:      y.Name,
			Number:    y.Number,
			Mintmark:  sqlOr(y.Mintmark, ""),
			Mintage:   sqlOr(y.Mintage, MintageUnknown),
			Reference: sqlOr(y.Reference, ""),
		})
	}

	return MintageData{xs, ys}, rs.Err()
}

func sqlOr[T any](v sql.Null[T], dflt T) T {
	if v.Valid {
		return v.V
	}
	return dflt
}

func denomToIdx(d float64) int {
	return slices.Index([]float64{
		0.01, 0.02, 0.05, 0.10,
		0.20, 0.50, 1.00, 2.00,
	}, d)
}
