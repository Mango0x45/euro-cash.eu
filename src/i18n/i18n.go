package i18n

import (
	"errors"
	"fmt"
	"maps"
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
	Bcp, Name                        string
	Eurozone, Enabled                bool
	DateFormat                       string
	GroupSeparator, DecimalSeparator rune
	MonetaryPre                      [2]string
	MonetaryPost                     string
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

	/* To determine the correct currency-, date-, and number formats to
	   use, use the ‘getfmt’ script in the repository root */
	locales = [...]LocaleInfo{
		{
			Bcp:              "ca",
			Name:             "Català",
			DateFormat:       "2/1/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "de",
			Name:             "Deutsch",
			DateFormat:       "2.1.2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "el",
			Name:             "Ελληνικά",
			DateFormat:       "2/1/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "en",
			Name:             "English",
			DateFormat:       "02/01/2006",
			Eurozone:         true,
			Enabled:          true,
			GroupSeparator:   ',',
			DecimalSeparator: '.',
			MonetaryPre:      [2]string{"€", "-€"},
		},
		{
			Bcp:              "es",
			Name:             "Español",
			DateFormat:       "2/1/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "et",
			Name:             "Eesti",
			DateFormat:       "2.1.2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "fi",
			Name:             "Suomi",
			DateFormat:       "2.1.2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "fr",
			Name:             "Français",
			DateFormat:       "02/01/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "ga",
			Name:             "Gaeilge",
			DateFormat:       "02/01/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ',',
			DecimalSeparator: '.',
			MonetaryPre:      [2]string{"€", "-€"},
		},
		{
			Bcp:              "hr",
			Name:             "Hrvatski",
			DateFormat:       "02. 01. 2006.",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "it",
			Name:             "Italiano",
			DateFormat:       "02/01/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "lb",
			Name:             "Lëtzebuergesch",
			DateFormat:       "2.1.2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "lt",
			Name:             "Lietuvių",
			DateFormat:       "2006-01-02",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "lv",
			Name:             "Latviešu",
			DateFormat:       "2.01.2006.",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "mt",
			Name:             "Malti",
			DateFormat:       "2/1/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ',',
			DecimalSeparator: '.',
			MonetaryPre:      [2]string{"€", "-€"},
		},
		{
			Bcp:              "nl",
			Name:             "Nederlands",
			DateFormat:       "2-1-2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"€ ", "€ -"},
		},
		{
			Bcp:              "pt",
			Name:             "Português",
			DateFormat:       "02/01/2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"€ ", "€ -"},
		},
		{
			Bcp:              "sk",
			Name:             "Slovenčina",
			DateFormat:       "2. 1. 2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "sl",
			Name:             "Slovenščina",
			DateFormat:       "2. 1. 2006",
			Eurozone:         true,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "sv",
			Name:             "Svenska",
			DateFormat:       "2006-01-02",
			Eurozone:         true,
			Enabled:          true,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		/* Non-Eurozone locales */
		{
			Bcp:              "bg",
			Name:             "Български",
			DateFormat:       "2.01.2006 г.",
			Eurozone:         false, /* TODO(2026): Set to true */
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "ro",
			Name:             "Română",
			DateFormat:       "02.01.2006",
			Eurozone:         false,
			Enabled:          false,
			GroupSeparator:   '.',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
		{
			Bcp:              "uk",
			Name:             "Yкраїнська",
			DateFormat:       "02.01.2006",
			Eurozone:         false,
			Enabled:          false,
			GroupSeparator:   ' ',
			DecimalSeparator: ',',
			MonetaryPre:      [2]string{"", "-"},
			MonetaryPost:     " €",
		},
	}
	Printers       map[string]Printer = make(map[string]Printer, len(locales))
	DefaultPrinter Printer
)

func Init() {
	for _, li := range locales {
		if !li.Enabled {
			continue
		}
		gl := gotext.NewLocale("po", li.Bcp)
		gl.AddDomain("messages")
		Printers[li.Bcp] = Printer{li, gl}
	}

	DefaultPrinter = Printers["en"]
}

func Locales() []LocaleInfo {
	return locales[:]
}

func (p Printer) Get(fmt string, args ...map[string]any) string {
	return p.Sprintf(p.inner.Get(fmt), args...)
}

func (p Printer) GetN(fmtS, fmtP string, n int, args ...map[string]any) string {
	return p.Sprintf(p.inner.GetN(fmtS, fmtP, n), args...)
}

/* Transform ‘en-US’ to ‘en’ */
func (l LocaleInfo) Language() string {
	return l.Bcp[:2]
}

func (p Printer) Itoa(n int) string {
	var bob strings.Builder
	writeInt(&bob, n, p.LocaleInfo)
	return bob.String()
}

func (p Printer) Ftoa(n float64) string {
	var bob strings.Builder
	writeFloat(&bob, n, p.LocaleInfo)
	return bob.String()
}

func (p Printer) Mitoa(n int) string {
	var bob strings.Builder
	sprintfm(p.LocaleInfo, &bob, n)
	return bob.String()
}

func (p Printer) Mftoa(n float64) string {
	var bob strings.Builder
	sprintfm(p.LocaleInfo, &bob, n)
	return bob.String()
}

func (p Printer) Sprintf(format string, args ...map[string]any) string {
	var bob strings.Builder
	vars := map[string]any{
		"-":    "a",
		"Null": "",
	}
	for _, arg := range args {
		maps.Copy(vars, arg)
	}

	for {
		i := strings.IndexByte(format, '{')
		if i == -1 {
			htmlesc(&bob, format)
			break
		}
		htmlesc(&bob, format[:i])

		format = format[i+1:]
		if len(format) == 0 {
			/* TODO: Handle error: trailing percent */
			break
		}

		i = strings.IndexRune(format, '}')
		if i == -1 {
			/* TODO: Handle error: unterminated { */
			return "unterminated {"
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

		v, ok := vars[parts[0]]
		if !ok {
			/* TODO: Handle error: no such key */
			return "no such key"
		}
		h(p.LocaleInfo, &bob, v)
	}

	return bob.String()
}

func sprintfGeneric(li LocaleInfo, bob *strings.Builder, v any) error {
	switch v.(type) {
	case time.Time:
		htmlesc(bob, v.(time.Time).Format(li.DateFormat))
	case int:
		writeInt(bob, v.(int), li)
	case float64:
		writeFloat(bob, v.(float64), li)
	case string:
		htmlesc(bob, v.(string))
	default:
		htmlesc(bob, fmt.Sprint(v))
	}
	return nil
}

func sprintfe(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("<a href=\"mailto:")
	htmlesc(bob, s)
	bob.WriteString("\">")
	htmlesc(bob, s)
	bob.WriteString("</a>")
	return nil
}

func sprintfE(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	for tag := range strings.SplitSeq(s, ",") {
		bob.WriteString("</")
		bob.WriteString(tag)
		bob.WriteByte('>')
	}
	return nil
}

func sprintfl(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("<a href=\"")
	htmlesc(bob, s)
	bob.WriteString("\">")
	return nil
}

func sprintfL(li LocaleInfo, bob *strings.Builder, v any) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("TODO")
	}
	bob.WriteString("<a href=\"")
	htmlesc(bob, s)
	bob.WriteString("\" target=\"_blank\">")
	return nil
}

func sprintfm(li LocaleInfo, bob *strings.Builder, v any) error {
	switch v.(type) {
	case int:
		n := v.(int)
		htmlesc(bob, li.MonetaryPre[btoi(n >= 0)])
		writeInt(bob, abs(n), li)
		htmlesc(bob, li.MonetaryPost)
	case float64:
		n := v.(float64)
		htmlesc(bob, li.MonetaryPre[btoi(n >= 0)])
		writeFloat(bob, abs(n), li)
		htmlesc(bob, li.MonetaryPost)
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
	bob.WriteString(s)
	return nil
}

func writeInt(bob *strings.Builder, num int, li LocaleInfo) {
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
			bob.WriteRune(li.GroupSeparator)
			c = 0
		}
	}
}

func writeFloat(bob *strings.Builder, num float64, li LocaleInfo) {
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
			bob.WriteRune(li.GroupSeparator)
			c = 0
		}
	}

	bob.WriteRune(li.DecimalSeparator)
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

func htmlesc(bob *strings.Builder, s string) {
	for _, r := range s {
		switch r {
		case '<':
			bob.WriteString("&lt;")
		case '>':
			bob.WriteString("&gt;")
		case '&':
			bob.WriteString("&amp;")
		case '"':
			bob.WriteString("&#34;")
		case '\'':
			bob.WriteString("&#39;")
		default:
			bob.WriteRune(r)
		}
	}
}

func htmlescByte(bob *strings.Builder, b byte) {
	switch b {
	case '<':
		bob.WriteString("&lt;")
	case '>':
		bob.WriteString("&gt;")
	case '&':
		bob.WriteString("&amp;")
	case '"':
		bob.WriteString("&#34;")
	case '\'':
		bob.WriteString("&#39;")
	default:
		bob.WriteByte(b)
	}
}
