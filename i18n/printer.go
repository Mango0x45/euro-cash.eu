package i18n

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:generate gotext -srclang=en-GB update -out=catalog.go -lang=en-GB,nl-NL git.thomasvoss.com/euro-cash.eu

type Printer struct {
	Locale  Locale
	printer *message.Printer
}

type Locale struct {
	Code     string
	Name     string
	dateFmt  string
	moneyFmt string
	Eurozone bool
	Enabled  bool
}

var (
	Locales = [...]Locale{
		Locale{
			Code:     "ca-AD",
			Name:     "català",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "de-DE",
			Name:     "Deutsch",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "el-GR",
			Name:     "ελληνικά",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "en-GB",
			Name:     "English",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  true,
		},
		Locale{
			Code:     "es-ES",
			Name:     "español",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "et-EE",
			Name:     "eesti",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "fi-FI",
			Name:     "suomi",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "fr-FR",
			Name:     "français",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "ga-IE",
			Name:     "Gaeilge",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "hr-HR",
			Name:     "hrvatski",
			dateFmt:  "02. 01. 2006.",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "it-IT",
			Name:     "italiano",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "lb-LU",
			Name:     "lëtzebuergesch",
			dateFmt:  "2.1.2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "lt-LT",
			Name:     "lietuvių",
			dateFmt:  "2006-01-02",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "lv-LV",
			Name:     "latviešu",
			dateFmt:  "2.01.2006.",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "mt-MT",
			Name:     "Malti",
			dateFmt:  "2/1/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "nl-NL",
			Name:     "Nederlands",
			dateFmt:  "2-1-2006",
			Eurozone: true,
			Enabled:  true,
		},
		Locale{
			Code:     "pt-PT",
			Name:     "português",
			dateFmt:  "02/01/2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "sk-SK",
			Name:     "slovenčina",
			dateFmt:  "2. 1. 2006",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "sl-SI",
			Name:     "slovenščina",
			dateFmt:  "2. 1. 2006",
			Eurozone: true,
			Enabled:  false,
		},

		/* Non-Eurozone locales */
		Locale{
			Code:     "bg-BG",
			Name:     "български",
			dateFmt:  "2.01.2006 г.",
			Eurozone: false,
			Enabled:  false,
		},
		Locale{
			Code:     "en-US",
			Name:     "English (US)",
			dateFmt:  "1/2/2006",
			Eurozone: false,
			Enabled:  false,
		},
		Locale{
			Code:     "ro-RO",
			Name:     "română",
			dateFmt:  "02.01.2006",
			Eurozone: false,
			Enabled:  false,
		},
		Locale{
			Code:     "uk-UA",
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
			lang := language.MustParse(loc.Code)
			Printers[strings.ToLower(loc.Code)] = Printer{
				Locale:  loc,
				printer: message.NewPrinter(lang),
			}
		}
	}
	DefaultPrinter = Printers["en-gb"]
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
	switch p.Locale.Code {
	case "en-GB", "en-US", "ga-IE", "mt-MT":
		return fmt.Sprintf("€%s", valstr)
	case "nl-NL":
		return fmt.Sprintf("€ %s", valstr)
	default:
		return fmt.Sprintf("%s €", valstr)
	}
}

func (l Locale) Language() string {
	return l.Code[:2]
}
