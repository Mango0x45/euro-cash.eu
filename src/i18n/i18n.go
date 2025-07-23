package i18n

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/leonelquinteros/gotext"
)

type Printer struct {
	LocaleInfo
	inner *gotext.Locale
}

type LocaleInfo struct {
	Bcp, Name                            string
	Eurozone, Enabled                    bool
	DateFormat                           string
	ThousandsSeparator, DecimalSeparator rune
	MonetaryPre, MonetaryPost            [2]string
}

type number interface {
	int | float64
}

type sprintfFunc func(LocaleInfo, *strings.Builder, any) error

var (
	handlers map[rune]sprintfFunc = map[rune]sprintfFunc{
		-1:  sprintfGeneric,
		'e': sprintfe,
		'E': sprintfE,
		'l': sprintfl,
		'L': sprintfL,
		'm': sprintfm,
		'r': sprintfr,
	}

	/* To determine the correct date format to use, use the ‘datefmt’ script in
	   the repository root */
	locales = [...]LocaleInfo{
		{
			Bcp:        "ca",
			Name:       "Català",
			DateFormat: "2/1/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "de",
			Name:       "Deutsch",
			DateFormat: "2.1.2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "el",
			Name:       "Ελληνικά",
			DateFormat: "2/1/2006",
			Eurozone:   true,
			Enabled:    true,
		},
		{
			Bcp:                "en",
			Name:               "English",
			DateFormat:         "02/01/2006",
			Eurozone:           true,
			Enabled:            true,
			ThousandsSeparator: ',',
			DecimalSeparator:   '.',
			MonetaryPre:        [2]string{"€", "-€"},
		},
		{
			Bcp:        "es",
			Name:       "Español",
			DateFormat: "2/1/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "et",
			Name:       "Eesti",
			DateFormat: "2.1.2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "fi",
			Name:       "Suomi",
			DateFormat: "2.1.2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "fr",
			Name:       "Français",
			DateFormat: "02/01/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "ga",
			Name:       "Gaeilge",
			DateFormat: "02/01/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "hr",
			Name:       "Hrvatski",
			DateFormat: "02. 01. 2006.",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "it",
			Name:       "Italiano",
			DateFormat: "02/01/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "lb",
			Name:       "Lëtzebuergesch",
			DateFormat: "2.1.2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "lt",
			Name:       "Lietuvių",
			DateFormat: "2006-01-02",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "lv",
			Name:       "Latviešu",
			DateFormat: "2.01.2006.",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "mt",
			Name:       "Malti",
			DateFormat: "2/1/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:                "nl",
			Name:               "Nederlands",
			DateFormat:         "2-1-2006",
			Eurozone:           true,
			Enabled:            true,
			ThousandsSeparator: '.',
			DecimalSeparator:   ',',
			MonetaryPre:        [2]string{"€ ", "€ -"},
		},
		{
			Bcp:        "pt",
			Name:       "Português",
			DateFormat: "02/01/2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "sk",
			Name:       "Slovenčina",
			DateFormat: "2. 1. 2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "sl",
			Name:       "Slovenščina",
			DateFormat: "2. 1. 2006",
			Eurozone:   true,
			Enabled:    false,
		},
		{
			Bcp:        "sv",
			Name:       "Svenska",
			DateFormat: "2006-01-02",
			Eurozone:   true,
			Enabled:    false,
		},

		/* Non-Eurozone locales */
		{
			Bcp:        "bg",
			Name:       "Български",
			DateFormat: "2.01.2006 г.",
			Eurozone:   false, /* TODO(2026): Set to true */
			Enabled:    true,
		},
		{
			Bcp:        "da",
			Name:       "Dansk",
			DateFormat: "02.01.2006",
			Eurozone:   false,
			Enabled:    false,
		},
		{
			Bcp:        "en-US",
			Name:       "English (US)",
			DateFormat: "1/2/2006",
			Eurozone:   false,
			Enabled:    false,
		},
		{
			Bcp:        "hu",
			Name:       "Magyar",
			DateFormat: "2006. 01. 02.",
			Eurozone:   false,
			Enabled:    false,
		},
		{
			Bcp:        "pl",
			Name:       "Polski",
			DateFormat: "2.01.2006",
			Eurozone:   false,
			Enabled:    false,
		},
		{
			Bcp:        "ro",
			Name:       "Română",
			DateFormat: "02.01.2006",
			Eurozone:   false,
			Enabled:    false,
		},
		{
			Bcp:        "uk",
			Name:       "Yкраїнська",
			DateFormat: "02.01.2006",
			Eurozone:   false,
			Enabled:    false,
		},
	}
	/* Map of language codes to printers.  We do this instead of just
	   using language.MustParse() directly so that we can easily see if a
	   language is supported or not. */
	Printers       map[string]Printer = make(map[string]Printer, len(locales))
	DefaultPrinter Printer
)

func Init() {
	for _, li := range locales {
		if !li.Enabled {
			continue
		}
		Printers[li.Bcp] = Printer{li, gotext.NewLocale("po", li.Bcp)}
	}

	DefaultPrinter = Printers["en"]
}

func Locales() []LocaleInfo {
	return locales[:]
}

func (p Printer) Get(fmt string, args ...map[string]any) string {
	/* TODO: Warning if you pass more than 1 arg? */
	var m map[string]any
	if len(args) == 0 {
		m = make(map[string]any)
	} else {
		m = args[0]
	}

	return p.Sprintf(p.inner.Get(fmt), m)
}

func (p Printer) GetN(fmtS, fmtP string, n int, args map[string]any) string {
	return p.Sprintf(p.inner.GetN(fmtS, fmtP, n), args)
}

/* Transform ‘en-US’ to ‘en’ */
func (l LocaleInfo) Language() string {
	return l.Bcp[:2]
}

func (p Printer) Sprintf(format string, args map[string]any) string {
	var bob strings.Builder
	args["-"] = ""

	for {
		i := strings.IndexByte(format, '%')
		if i == -1 {
			bob.WriteString(format)
			break
		}
		bob.WriteString(format[:i])

		format = format[i+1:]
		if len(format) == 0 {
			/* TODO: Handle error: trailing percent */
			break
		}

		b := format[0]
		format = format[1:]

		switch b {
		case '%':
			bob.WriteByte(b)
		case '(':
			i = strings.IndexRune(format, ')')
			if i == -1 {
				/* TODO: Handle error: unterminated %( */
				return "unterminated %("
			}

			parts := strings.Split(format[:i], ":")
			format = format[i+1:]

			var flag rune
			switch len(parts) {
			case 1:
				flag = -1
			case 2:
				f, n := utf8.DecodeRune([]byte(parts[1]))
				if n != len(parts[1]) {
					/* TODO: Handle error: flag too long or empty */
					return "flag too long or empty"
				}
				flag = f
			default:
				/* TODO: Handle error: too many colons */
				return "too many colons"
			}

			h, ok := handlers[flag]
			if !ok {
				/* TODO: Handle error: no such handler */
				return "no such handler"
			}

			v, ok := args[parts[0]]
			if !ok {
				/* TODO: Handle error: no such key */
				return "no such key"
			}
			h(p.LocaleInfo, &bob, v)
		default:
			/* TODO: Handle error: invalid escape */
			bob.WriteByte(b)
		}
	}

	return bob.String()
}

func sprintfGeneric(li LocaleInfo, bob *strings.Builder, v any) error {
	switch v.(type) {
	case time.Time:
		bob.WriteString(v.(time.Time).Format(li.DateFormat))
	case int:
		writeInt(bob, v.(int), li.ThousandsSeparator)
	case float64:
		writeFloat(bob, v.(float64), li.ThousandsSeparator, li.DecimalSeparator)
	case string:
		bob.WriteString(v.(string))
	default:
		bob.WriteString(fmt.Sprint(v))
	}
	return nil
}

func sprintfe(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("\u2068<a href=\"mailto:")
	attrEscape(bob, s)
	bob.WriteString("\">\u2069")
	bob.WriteString(s)
	bob.WriteString("\u2068</a>\u2069")
	return nil
}

func sprintfE(li LocaleInfo, bob *strings.Builder, _ any) error {
	bob.WriteString("\u2068</a>\u2069")
	return nil
}

func sprintfl(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("\u2068<a href=\"")
	attrEscape(bob, s)
	bob.WriteString("\">\u2069")
	return nil
}

func sprintfL(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("\u2068<a href=\"")
	attrEscape(bob, s)
	bob.WriteString("\" target=\"_blank\">\u2069")
	return nil
}

func sprintfm(li LocaleInfo, bob *strings.Builder, v any) error {
	switch v.(type) {
	case int:
		n := v.(int)
		i := btoi(n >= 0)
		bob.WriteString(li.MonetaryPre[i])
		writeInt(bob, abs(n), li.ThousandsSeparator)
		bob.WriteString(li.MonetaryPost[i])
	case float64:
		n := v.(float64)
		i := btoi(n >= 0)
		bob.WriteString(li.MonetaryPre[i])
		writeFloat(bob, abs(n), li.ThousandsSeparator, li.DecimalSeparator)
		bob.WriteString(li.MonetaryPost[i])
	default:
		return errors.New("TODO")
	}
	return nil
}

func sprintfr(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteRune('\u2068')
	bob.WriteString(s)
	bob.WriteRune('\u2069')
	return nil
}

func attrEscape(bob *strings.Builder, s string) {
	for _, r := range s {
		switch r {
		case '<':
			bob.WriteString("&lt;")
		case '&':
			bob.WriteString("&amp;")
		case '"':
			bob.WriteString("&quot;")
		default:
			bob.WriteRune(r)
		}
	}
}

func writeInt(bob *strings.Builder, num int, sep rune) {
	s := fmt.Sprintf("%d", num)
	if s[0] == '-' {
		bob.WriteByte('-')
		s = s[1:]
	}
	n := len(s)
	c := 3 - n%3
	if c == 3 {
		c = 0
	}
	for i := 0; i < n; i++ {
		c++
		bob.WriteByte(s[i])
		if c == 3 && i+1 < n {
			bob.WriteRune(sep)
			c = 0
		}
	}
}

func writeFloat(bob *strings.Builder, num float64, tsep, dsep rune) {
	s := fmt.Sprintf("%.2f", num)
	if s[0] == '-' {
		bob.WriteByte('-')
		s = s[1:]
	}

	n := strings.IndexByte(s, '.')
	c := 3 - n%3
	if c == 3 {
		c = 0
	}
	for i := 0; i < n; i++ {
		c++
		bob.WriteByte(s[i])
		if c == 3 && i+1 < n {
			bob.WriteRune(tsep)
			c = 0
		}
	}

	bob.WriteRune(dsep)
	bob.WriteString(s[n+1:])
}

func abs[T number](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func btoi(b bool) int {
	if b {
		return 0
	}
	return 1
}
