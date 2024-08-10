//go:generate templ generate -log-level warn

package templates

import "git.thomasvoss.com/euro-cash.eu/i18n"

type country struct{ code, name string }

func countries(p i18n.Printer) []country {
	return []country{
		{code: "AD", name: p.T("Andorra")},
		{code: "AT", name: p.T("Austria")},
		{code: "BE", name: p.T("Belgium")},
		{code: "CY", name: p.T("Cyprus")},
		{code: "DE", name: p.T("Germany")},
		{code: "EE", name: p.T("Estonia")},
		{code: "ES", name: p.T("Spain")},
		{code: "FI", name: p.T("Finland")},
		{code: "FR", name: p.T("France")},
		{code: "GR", name: p.T("Greece")},
		{code: "HR", name: p.T("Croatia")},
		{code: "IE", name: p.T("Ireland")},
		{code: "IT", name: p.T("Italy")},
		{code: "LT", name: p.T("Lithuania")},
		{code: "LU", name: p.T("Luxembourg")},
		{code: "LV", name: p.T("Latvia")},
		{code: "MC", name: p.T("Monaco")},
		{code: "MT", name: p.T("Malta")},
		{code: "NL", name: p.T("Netherlands")},
		{code: "PT", name: p.T("Portugal")},
		{code: "SI", name: p.T("Slovenia")},
		{code: "SK", name: p.T("Slovakia")},
		{code: "SM", name: p.T("San Marino")},
		{code: "VA", name: p.T("Vatican City")},
	}
}
