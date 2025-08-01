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

func sortedCountries(p i18n.Printer) []country {
	xs := []country{
		{Code: "ad", Name: p.GetC("Andorra", "Country Name")},
		{Code: "at", Name: p.GetC("Austria", "Country Name")},
		{Code: "be", Name: p.GetC("Belgium", "Country Name")},
		/* TODO(2026): Add Bulgaria */
		/* {Code: "bg", Name: p.GetC("Bulgaria", "Country Name")}, */
		{Code: "cy", Name: p.GetC("Cyprus", "Country Name")},
		{Code: "de", Name: p.GetC("Germany", "Country Name")},
		{Code: "ee", Name: p.GetC("Estonia", "Country Name")},
		{Code: "es", Name: p.GetC("Spain", "Country Name")},
		{Code: "fi", Name: p.GetC("Finland", "Country Name")},
		{Code: "fr", Name: p.GetC("France", "Country Name")},
		{Code: "gr", Name: p.GetC("Greece", "Country Name")},
		{Code: "hr", Name: p.GetC("Croatia", "Country Name")},
		{Code: "ie", Name: p.GetC("Ireland", "Country Name")},
		{Code: "it", Name: p.GetC("Italy", "Country Name")},
		{Code: "lt", Name: p.GetC("Lithuania", "Country Name")},
		{Code: "lu", Name: p.GetC("Luxembourg", "Country Name")},
		{Code: "lv", Name: p.GetC("Latvia", "Country Name")},
		{Code: "mc", Name: p.GetC("Monaco", "Country Name")},
		{Code: "mt", Name: p.GetC("Malta", "Country Name")},
		{Code: "nl", Name: p.GetC("Netherlands", "Country Name")},
		{Code: "pt", Name: p.GetC("Portugal", "Country Name")},
		{Code: "si", Name: p.GetC("Slovenia", "Country Name")},
		{Code: "sk", Name: p.GetC("Slovakia", "Country Name")},
		{Code: "sm", Name: p.GetC("San Marino", "Country Name")},
		{Code: "va", Name: p.GetC("Vatican City", "Country Name")},
	}
	c := collate.New(language.MustParse(p.Bcp))
	slices.SortFunc(xs, func(x, y country) int {
		return c.CompareString(x.Name, y.Name)
	})
	return xs
}
