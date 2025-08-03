package app

import (
	"slices"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

type country struct {
	Code, Name string
}

var countryCodeToName = map[string]string{
	"ad": "Andorra",
	"at": "Austria",
	"be": "Belgium",
	/* TODO(2026): Add Bulgaria */
	/* "bg": "Bulgaria", */
	"cy": "Cyprus",
	"de": "Germany",
	"ee": "Estonia",
	"es": "Spain",
	"fi": "Finland",
	"fr": "France",
	"gr": "Greece",
	"hr": "Croatia",
	"ie": "Ireland",
	"it": "Italy",
	"lt": "Lithuania",
	"lu": "Luxembourg",
	"lv": "Latvia",
	"mc": "Monaco",
	"mt": "Malta",
	"nl": "Netherlands",
	"pt": "Portugal",
	"si": "Slovenia",
	"sk": "Slovakia",
	"sm": "San Marino",
	"va": "Vatican City",
}

func sortedCountries(p i18n.Printer) []country {
	xs := []country{
		{Code: "ad", Name: p.GetC("Andorra", "Place Name")},
		{Code: "at", Name: p.GetC("Austria", "Place Name")},
		{Code: "be", Name: p.GetC("Belgium", "Place Name")},
		/* TODO(2026): Add Bulgaria */
		/* {Code: "bg", Name: p.GetC("Bulgaria", "Place Name")}, */
		{Code: "cy", Name: p.GetC("Cyprus", "Place Name")},
		{Code: "de", Name: p.GetC("Germany", "Place Name")},
		{Code: "ee", Name: p.GetC("Estonia", "Place Name")},
		{Code: "es", Name: p.GetC("Spain", "Place Name")},
		{Code: "fi", Name: p.GetC("Finland", "Place Name")},
		{Code: "fr", Name: p.GetC("France", "Place Name")},
		{Code: "gr", Name: p.GetC("Greece", "Place Name")},
		{Code: "hr", Name: p.GetC("Croatia", "Place Name")},
		{Code: "ie", Name: p.GetC("Ireland", "Place Name")},
		{Code: "it", Name: p.GetC("Italy", "Place Name")},
		{Code: "lt", Name: p.GetC("Lithuania", "Place Name")},
		{Code: "lu", Name: p.GetC("Luxembourg", "Place Name")},
		{Code: "lv", Name: p.GetC("Latvia", "Place Name")},
		{Code: "mc", Name: p.GetC("Monaco", "Place Name")},
		{Code: "mt", Name: p.GetC("Malta", "Place Name")},
		{Code: "nl", Name: p.GetC("Netherlands", "Place Name")},
		{Code: "pt", Name: p.GetC("Portugal", "Place Name")},
		{Code: "si", Name: p.GetC("Slovenia", "Place Name")},
		{Code: "sk", Name: p.GetC("Slovakia", "Place Name")},
		{Code: "sm", Name: p.GetC("San Marino", "Place Name")},
		{Code: "va", Name: p.GetC("Vatican City", "Place Name")},
	}
	c := collate.New(language.MustParse(p.Bcp))
	slices.SortFunc(xs, func(x, y country) int {
		return c.CompareString(x.Name, y.Name)
	})
	return xs
}
