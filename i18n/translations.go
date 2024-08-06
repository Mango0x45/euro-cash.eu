package i18n

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:generate gotext -srclang=en-GB update -out=catalog.go -lang=en-GB,nl-NL git.thomasvoss.com/euro-cash.eu

type Printer struct {
	Lang    string
	printer *message.Printer
}

type Locale struct {
	Code     string
	Name     string
	Eurozone bool
	Enabled  bool
}

var (
	Locales = [...]Locale{
		Locale{
			Code:     "ca-AD",
			Name:     "català",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "de-DE",
			Name:     "Deutsch",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "el-GR",
			Name:     "ελληνικά",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "en-GB",
			Name:     "English",
			Eurozone: true,
			Enabled:  true,
		},
		Locale{
			Code:     "es-ES",
			Name:     "español",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "et-EE",
			Name:     "eesti",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "fi-FI",
			Name:     "suomi",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "fr-FR",
			Name:     "français",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "ga-IE",
			Name:     "Gaeilge",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "hr-HR",
			Name:     "hrvatski",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "it-IT",
			Name:     "italiano",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "lt-LT",
			Name:     "lietuvių",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "lv-LV",
			Name:     "latviešu",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "mt-MT",
			Name:     "Malti",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "nl-NL",
			Name:     "Nederlands",
			Eurozone: true,
			Enabled:  true,
		},
		Locale{
			Code:     "pt-PT",
			Name:     "português",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "sk-SK",
			Name:     "slovenčina",
			Eurozone: true,
			Enabled:  false,
		},
		Locale{
			Code:     "sl-SI",
			Name:     "slovenščina",
			Eurozone: true,
			Enabled:  false,
		},

		/* Non-Eurozone locales */
		Locale{
			Code:     "bg-BG",
			Name:     "български",
			Eurozone: false,
			Enabled:  false,
		},
		Locale{
			Code:     "en-US",
			Name:     "English (US)",
			Eurozone: false,
			Enabled:  false,
		},
		Locale{
			Code:     "ro-RO",
			Name:     "română",
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
				Lang:    loc.Code,
				printer: message.NewPrinter(lang),
			}
		}
	}
	DefaultPrinter = Printers["en-gb"]
}

func (p Printer) T(fmt string, args ...any) string {
	return p.printer.Sprintf(fmt, args...)
}

func (l Locale) Language() string {
	return l.Code[:2]
}
