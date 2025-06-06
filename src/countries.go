package src

import (
	"slices"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type country struct {
	Code, Name string
}

func sortedCountries(p Printer) []country {
	xs := []country{
		{Code: "ad", Name: p.T("Andorra")},
		{Code: "at", Name: p.T("Austria")},
		{Code: "be", Name: p.T("Belgium")},
		/* TODO(2026): Add Bulgaria */
		/* {Code: "bg", Name: p.T("Bulgaria")}, */
		{Code: "cy", Name: p.T("Cyprus")},
		{Code: "de", Name: p.T("Germany")},
		{Code: "ee", Name: p.T("Estonia")},
		{Code: "es", Name: p.T("Spain")},
		{Code: "fi", Name: p.T("Finland")},
		{Code: "fr", Name: p.T("France")},
		{Code: "gr", Name: p.T("Greece")},
		{Code: "hr", Name: p.T("Croatia")},
		{Code: "ie", Name: p.T("Ireland")},
		{Code: "it", Name: p.T("Italy")},
		{Code: "lt", Name: p.T("Lithuania")},
		{Code: "lu", Name: p.T("Luxembourg")},
		{Code: "lv", Name: p.T("Latvia")},
		{Code: "mc", Name: p.T("Monaco")},
		{Code: "mt", Name: p.T("Malta")},
		{Code: "nl", Name: p.T("Netherlands")},
		{Code: "pt", Name: p.T("Portugal")},
		{Code: "si", Name: p.T("Slovenia")},
		{Code: "sk", Name: p.T("Slovakia")},
		{Code: "sm", Name: p.T("San Marino")},
		{Code: "va", Name: p.T("Vatican City")},
	}
	c := collate.New(language.MustParse(p.Locale.Bcp))
	slices.SortFunc(xs, func(x, y country) int {
		return c.CompareString(x.Name, y.Name)
	})
	return xs
}
