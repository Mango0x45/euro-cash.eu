//go:generate go tool gotext -srclang=en -dir=rosetta extract -lang=bg,el,en,nl .
//go:generate ../exttmpl

package src

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"git.thomasvoss.com/euro-cash.eu/src/email"
)

type Printer struct {
	Locale locale
	inner  *message.Printer
}

type locale struct {
	Bcp, Name         string
	Eurozone, Enabled bool
	dateFmt           string
}

var (
    /* To determine the correct date format to use, use the ‘datefmt’ script in
       the repository root */
	Locales = [...]locale{
		{
			Bcp:      "ca",
			Name:     "català",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "de",
			Name:     "Deutsch",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "el",
			Name:     "ελληνικά",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  true,
		},
		{
			Bcp:      "en",
			Name:     "English",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  true,
		},
		{
			Bcp:      "es",
			Name:     "español",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "et",
			Name:     "eesti",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "fi",
			Name:     "suomi",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "fr",
			Name:     "français",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "ga",
			Name:     "Gaeilge",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "hr",
			Name:     "hrvatski",
			dateFmt:  "02. 01. 2006.",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "it",
			Name:     "italiano",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "lb",
			Name:     "lëtzebuergesch",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "lt",
			Name:     "lietuvių",
			dateFmt:  "2006-01-02",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "lv",
			Name:     "latviešu",
			dateFmt:  "2.01.2006.",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "mt",
			Name:     "Malti",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "nl",
			Name:     "Nederlands",
			dateFmt:  "2-1-2006",
			Eurozone: true,
			Enabled:  true,
		},
		{
			Bcp:      "pt",
			Name:     "português",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "sk",
			Name:     "slovenčina",
			dateFmt:  "2. 1. 2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "sl",
			Name:     "slovenščina",
			dateFmt:  "2. 1. 2006",
			Eurozone: true,
			Enabled:  false,
		},
		{
			Bcp:      "sv",
			Name:     "svenska",
			dateFmt:  "2006-01-02",
			Eurozone: true,
			Enabled:  false,
		},

		/* Non-Eurozone locales */
		{
			Bcp:      "bg",
			Name:     "български",
			dateFmt:  "2.01.2006 г.",
			Eurozone: false,    /* TODO(2026): Set to true */
			Enabled:  true,
		},
		{
			Bcp:      "da",
			Name:     "dansk",
			dateFmt:  "02.01.2006",
			Eurozone: false,
			Enabled:  false,
		},
		{
			Bcp:      "en-US",
			Name:     "English (US)",
			dateFmt:  "1/2/2006",
			Eurozone: false,
			Enabled:  false,
		},
		{
			Bcp:      "hu",
			Name:     "magyar",
			dateFmt:  "2006. 01. 02.",
			Eurozone: false,
			Enabled:  false,
		},
		{
			Bcp:      "pl",
			Name:     "polski",
			dateFmt:  "2.01.2006",
			Eurozone: false,
			Enabled:  false,
		},
		{
			Bcp:      "ro",
			Name:     "română",
			dateFmt:  "02.01.2006",
			Eurozone: false,
			Enabled:  false,
		},
		{
			Bcp:      "uk",
			Name:     "yкраїнська",
			dateFmt:  "02.01.2006",
			Eurozone: false,
			Enabled:  false,
		},
	}
	/* Map of language codes to printers.  We do this instead of just
	   using language.MustParse() directly so that we can easily see if a
	   language is supported or not. */
	printers       map[string]Printer = make(map[string]Printer, len(Locales))
	defaultPrinter Printer
)

func init() {
	for _, loc := range Locales {
		if loc.Enabled {
			lang := language.MustParse(loc.Bcp)
			printers[strings.ToLower(loc.Bcp)] = Printer{
				Locale: loc,
				inner:  message.NewPrinter(lang),
			}
		}
	}
	defaultPrinter = printers["en"]
}

func (p Printer) T(fmt string, args ...any) string {
	return p.inner.Sprintf(fmt, args...)
}

func (p Printer) N(n int) string {
	return p.inner.Sprint(n)
}

func (p Printer) Date(d time.Time) string {
	return d.Format(p.Locale.dateFmt)
}

func (p Printer) M(val any) string {
	var vstr string

	/* Hack to avoid gotext writing these two ‘translations’ into the
	   translations file */
	f := p.inner.Sprintf

	switch val.(type) {
	case int:
		vstr = f("%d", val)
	case float64:
		vstr = f("%.2f", val)
	default:
		go func() {
			if err := email.ServerError(badMType{"TODO"}); err != nil {
				log.Print(err)
			}
		}()
		/* Hopefully this never happens */
		vstr = "ERROR"
	}

	/* All Eurozone languages place the eurosign after the value except
	   for Dutch, English, Irish, and Maltese.  Austrian German also
	   uses Dutch-style formatting, but we do not support that dialect. */
	switch p.Locale.Bcp {
	case "en", "en-US", "ga", "mt":
		return fmt.Sprintf("€%s", vstr)
	case "nl":
		return fmt.Sprintf("€ %s", vstr)
	default:
		return fmt.Sprintf("%s €", vstr)
	}
}

/* Transform ‘en-US’ to ‘en’ */
func (l locale) Language() string {
	return l.Bcp[:2]
}

type badMType struct {
	inner any
}

func (e badMType) Error() string {
	return fmt.Sprintf(
		"Attempted to format ‘%v’ of type ‘%T’ as a monetary value in Printer.M()",
		e.inner, e.inner)
}
