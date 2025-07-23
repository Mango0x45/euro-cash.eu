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
		{Code: "ad", Name: p.Get("Andorra")},
		{Code: "at", Name: p.Get("Austria")},
		{Code: "be", Name: p.Get("Belgium")},
		/* TODO(2026): Add Bulgaria */
		/* {Code: "bg", Name: p.Get("Bulgaria")}, */
		{Code: "cy", Name: p.Get("Cyprus")},
		{Code: "de", Name: p.Get("Germany")},
		{Code: "ee", Name: p.Get("Estonia")},
		{Code: "es", Name: p.Get("Spain")},
		{Code: "fi", Name: p.Get("Finland")},
		{Code: "fr", Name: p.Get("France")},
		{Code: "gr", Name: p.Get("Greece")},
		{Code: "hr", Name: p.Get("Croatia")},
		{Code: "ie", Name: p.Get("Ireland")},
		{Code: "it", Name: p.Get("Italy")},
		{Code: "lt", Name: p.Get("Lithuania")},
		{Code: "lu", Name: p.Get("Luxembourg")},
		{Code: "lv", Name: p.Get("Latvia")},
		{Code: "mc", Name: p.Get("Monaco")},
		{Code: "mt", Name: p.Get("Malta")},
		{Code: "nl", Name: p.Get("Netherlands")},
		{Code: "pt", Name: p.Get("Portugal")},
		{Code: "si", Name: p.Get("Slovenia")},
		{Code: "sk", Name: p.Get("Slovakia")},
		{Code: "sm", Name: p.Get("San Marino")},
		{Code: "va", Name: p.Get("Vatican City")},
	}
	c := collate.New(language.MustParse(p.Bcp))
	slices.SortFunc(xs, func(x, y country) int {
		return c.CompareString(x.Name, y.Name)
	})
	return xs
}
