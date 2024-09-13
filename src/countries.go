package src

import (
	"slices"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type country struct {
	code, name string
}

func sortedCountries(p Printer) []country {
	xs := []country{
		{code: "ad", name: p.T("Andorra")},
		{code: "at", name: p.T("Austria")},
		{code: "be", name: p.T("Belgium")},
		{code: "cy", name: p.T("Cyprus")},
		{code: "de", name: p.T("Germany")},
		{code: "ee", name: p.T("Estonia")},
		{code: "es", name: p.T("Spain")},
		{code: "fi", name: p.T("Finland")},
		{code: "fr", name: p.T("France")},
		{code: "gr", name: p.T("Greece")},
		{code: "hr", name: p.T("Croatia")},
		{code: "ie", name: p.T("Ireland")},
		{code: "it", name: p.T("Italy")},
		{code: "lt", name: p.T("Lithuania")},
		{code: "lu", name: p.T("Luxembourg")},
		{code: "lv", name: p.T("Latvia")},
		{code: "mc", name: p.T("Monaco")},
		{code: "mt", name: p.T("Malta")},
		{code: "nl", name: p.T("Netherlands")},
		{code: "pt", name: p.T("Portugal")},
		{code: "si", name: p.T("Slovenia")},
		{code: "sk", name: p.T("Slovakia")},
		{code: "sm", name: p.T("San Marino")},
		{code: "va", name: p.T("Vatican City")},
	}
	c := collate.New(language.MustParse(p.Locale.Bcp))
	slices.SortFunc(xs, func(x, y country) int {
		return c.CompareString(x.name, y.name)
	})
	return xs
}
