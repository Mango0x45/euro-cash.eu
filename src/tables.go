package app

import (
	"database/sql"
	"slices"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"git.thomasvoss.com/euro-cash.eu/src/dbx"
	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

type YearCCsPair struct {
	Year     int
	Mintmark sql.Null[string]
	CCs      []dbx.MCommemorative
}

type CountryMintageTable struct {
	Standard      []dbx.MSCountryRow
	Commemorative []YearCCsPair
}

type CountryCCsPair struct {
	Country  string
	Mintmark sql.Null[string]
	CCs      []dbx.MCommemorative
}

type YearMintageTable struct {
	Standard      []dbx.MSYearRow
	Commemorative []CountryCCsPair
}

func makeCountryMintageTable(data dbx.CountryMintageData, p i18n.Printer) CountryMintageTable {
	ccdata := data.Commemorative
	ccs := make([]YearCCsPair, 0, len(ccdata))

	for len(ccdata) > 0 {
		x := ccdata[0]
		i := len(ccdata)
		for j, y := range ccdata[1:] {
			if x.Year != y.Year || x.Mintmark != y.Mintmark {
				i = j + 1
				break
			}
		}
		ccs = append(ccs, YearCCsPair{
			Year:     x.Year,
			Mintmark: x.Mintmark,
			CCs:      ccdata[:i],
		})
		ccdata = ccdata[i:]
	}

	return CountryMintageTable{data.Standard, ccs}
}

func makeYearMintageTable(data dbx.YearMintageData, p i18n.Printer) YearMintageTable {
	ccdata := data.Commemorative
	ccs := make([]CountryCCsPair, 0, len(ccdata))

	for len(ccdata) > 0 {
		x := ccdata[0]
		i := len(ccdata)
		for j, y := range ccdata[1:] {
			if x.Country != y.Country || x.Mintmark != y.Mintmark {
				i = j + 1
				break
			}
		}
		ccs = append(ccs, CountryCCsPair{
			Country:  x.Country,
			Mintmark: x.Mintmark,
			CCs:      ccdata[:i],
		})
		ccdata = ccdata[i:]
	}

	/* NOTE: It’s safe to use MustParse() here, because by this
	   point we know that all BCPs are valid. */
	c := collate.New(language.MustParse(p.Bcp))
	for i, r := range data.Standard {
		name := countryCodeToName[r.Country]
		data.Standard[i].Country = p.GetC(name, "Place Name")
	}
	for i, r := range ccs {
		name := countryCodeToName[r.Country]
		ccs[i].Country = p.GetC(name, "Place Name")
	}
	slices.SortFunc(data.Standard, func(x, y dbx.MSYearRow) int {
		Δ := c.CompareString(x.Country, y.Country)
		if Δ == 0 {
			Δ = c.CompareString(x.Mintmark.V, y.Mintmark.V)
		}
		return Δ
	})
	slices.SortFunc(ccs, func(x, y CountryCCsPair) int {
		Δ := c.CompareString(x.Country, y.Country)
		if Δ == 0 {
			Δ = c.CompareString(x.Mintmark.V, y.Mintmark.V)
		}
		return Δ
	})

	return YearMintageTable{data.Standard, ccs}
}
