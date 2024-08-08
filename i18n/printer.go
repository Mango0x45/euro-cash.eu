package i18n

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:generate gotext -srclang=en update -out=catalog.go -lang=el,en,nl git.thomasvoss.com/euro-cash.eu

type Printer struct {
	Locale  Locale
	printer *message.Printer
}

type Locale struct {
	Bcp      string
	Name     string
	dateFmt  string
	moneyFmt string
	Eurozone bool
	Enabled  bool
}

var (
	Locales = [...]Locale{
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

		/* Non-Eurozone locales */
		{
			Bcp:      "bg",
			Name:     "български",
			dateFmt:  "2.01.2006 г.",
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
	Printers       map[string]Printer = make(map[string]Printer, len(Locales))
	DefaultPrinter Printer
)

func InitPrinters() {
	for _, loc := range Locales {
		if loc.Enabled {
			lang := language.MustParse(loc.Bcp)
			Printers[strings.ToLower(loc.Bcp)] = Printer{
				Locale:  loc,
				printer: message.NewPrinter(lang),
			}
		}
	}
	DefaultPrinter = Printers["en"]
}

func (p Printer) T(fmt string, args ...any) string {
	return p.printer.Sprintf(fmt, args...)
}

func (p Printer) Date(d time.Time) string {
	return d.Format(p.Locale.dateFmt)
}

/* TODO: Try to use a decimal type here */
func (p Printer) Money(val float64, round bool) string {
	var valstr string

	/* Hack to avoid gotext writing these two ‘translations’ into the
	   translations file */
	f := p.printer.Sprintf
	if round {
		valstr = f("%d", int(val))
	} else {
		valstr = f("%.2f", val)
	}

	/* All Eurozone languages place the eurosign after the value except
	   for Dutch, English, Gaelic, and Maltese.  Austrian German also
	   uses Dutch-style formatting, but we do not support that dialect. */
	switch p.Locale.Bcp {
	case "en", "en-US", "ga", "mt":
		return fmt.Sprintf("€%s", valstr)
	case "nl":
		return fmt.Sprintf("€ %s", valstr)
	default:
		return fmt.Sprintf("%s €", valstr)
	}
}

/* Transform ‘en-US’ to ‘en’ */
func (l Locale) Language() string {
	return l.Bcp[:2]
}
