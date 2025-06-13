package dbx

type MintageData struct {
	Standard      []MSRow
	Commemorative []MCRow
}

type MSRow struct {
	Type      int          `db:"type"`
	Year      int          `db:"year"`
	Mintmark  string       `db:"mintmark"`
	Mintages  [ndenoms]int `db:"€0,01;€0,02;€0,05;€0,10;€0,20;€0,50;€1,00;€2,00"`
	Reference string       `db:"reference"`
}

type MCRow struct {
	Type      int    `db:"type"`
	Year      int    `db:"year"`
	Name      string `db:"name"`
	Number    int    `db:"number"`
	Mintmark  string `db:"mintmark"`
	Mintage   int    `db:"mintage"`
	Reference string `db:"reference"`
}

/* DO NOT REORDER! */
const (
	TypeCirc = iota
	TypeNifc
	TypeProof
)

/* DO NOT REORDER! */
const (
	MintageUnknown = -iota - 1
	MintageInvalid
)

const ndenoms = 8

func GetMintages(country string) (MintageData, error) {
	var zero MintageData

	srows, err := db.Query(`SELECT * FROM mintages_s WHERE country = ?`, country)
	if err != nil {
		return zero, err
	}
	defer srows.Close()
	xs, err := scanToStructs[MSRow](srows)
	if err != nil {
		return zero, err
	}

	crows, err := db.Query(`SELECT * FROM mintages_c WHERE country = ?`, country)
	if err != nil {
		return zero, err
	}
	defer crows.Close()
	ys, err := scanToStructs[MCRow](crows)
	if err != nil {
		return zero, err
	}

	return MintageData{xs, ys}, nil
}
