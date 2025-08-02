//go:generate ./gen.py

package i18n

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"git.thomasvoss.com/euro-cash.eu/pkg/atexit"
	"git.thomasvoss.com/euro-cash.eu/pkg/watch"
	"github.com/leonelquinteros/gotext"
)

type Printer struct {
	LocaleInfo
	inner *gotext.Locale
}

type number interface {
	int | float64
}

type sprintfFunc func(LocaleInfo, *strings.Builder, any) error

var (
	handlers map[rune]sprintfFunc = map[rune]sprintfFunc{
		-1:  sprintfGeneric,
		'E': sprintfE,
		'L': sprintfL,
		'e': sprintfe,
		'l': sprintfl,
		'm': sprintfm,
		'p': sprintfp,
		'r': sprintfr,
	}

	Printers       map[string]Printer = make(map[string]Printer, len(locales))
	DefaultPrinter Printer
)

func Init(dir fs.FS, debugp bool) {
	gotext.FallbackLocale = "en"
	i := slices.IndexFunc(locales[:], func(li LocaleInfo) bool {
		return li.Bcp == gotext.FallbackLocale
	})
	if i == -1 {
		atexit.Exec()
		log.Fatalf("No translation file default locale ‘%s’\n",
			gotext.FallbackLocale)
	}
	if !locales[i].Enabled {
		atexit.Exec()
		log.Fatalf("Default locale ‘%s’ is not enabled\n",
			locales[i].Name)
	}

	initLocale(dir, locales[i], locales[i].Name, debugp)
	DefaultPrinter = Printers[gotext.FallbackLocale]

	for j, li := range locales {
		if li.Enabled && i != j {
			name := DefaultPrinter.GetC(li.Name, "Language Name")
			initLocale(dir, li, name, debugp)
		}
	}
}

func initLocale(dir fs.FS, li LocaleInfo, name string, debugp bool) {
	gl := gotext.NewLocaleFS(li.Bcp, dir)
	gl.AddDomain("messages")
	Printers[li.Bcp] = Printer{li, gl}

	if debugp {
		subdir, err := fs.Sub(dir, li.Bcp)
		if err != nil {
			log.Printf("No translations directory for ‘%s’\n", name)
			return
		}
		go watch.FileFS(subdir, "messages.po", func() {
			Printers[li.Bcp].inner.AddDomain("messages")
			log.Printf("Translations for ‘%s’ updated\n", name)
		})
	}

	log.Printf("Initialized printer for ‘%s’\n", name)
}

func Locales() []LocaleInfo {
	return locales[:]
}

func (p Printer) Get(fmt string, args ...map[string]any) string {
	return p.Sprintf(p.inner.Get(fmt), args...)
}

func (p Printer) GetC(fmt, ctx string, args ...map[string]any) string {
	return p.Sprintf(p.inner.GetC(fmt, ctx), args...)
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

func (p Printer) Itop(n int) string {
	var bob strings.Builder
	sprintfp(p.LocaleInfo, &bob, n)
	return bob.String()
}

func (p Printer) Ftop(n float64) string {
	var bob strings.Builder
	sprintfp(p.LocaleInfo, &bob, n)
	return bob.String()
}

func (p Printer) Itom(n int) string {
	var bob strings.Builder
	sprintfm(p.LocaleInfo, &bob, n)
	return bob.String()
}

func (p Printer) Ftom(n float64) string {
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
		htmlesc(bob, li.MonetarySuf)
	case float64:
		n := v.(float64)
		htmlesc(bob, li.MonetaryPre[btoi(n >= 0)])
		writeFloat(bob, abs(n), li)
		htmlesc(bob, li.MonetarySuf)
	default:
		return errors.New("TODO")
	}
	return nil
}

func sprintfp(li LocaleInfo, bob *strings.Builder, v any) error {
	var bob2 strings.Builder
	switch v.(type) {
	case int:
		writeInt(&bob2, v.(int), li)
	case float64:
		writeFloat(&bob2, v.(float64), li)
	default:
		return errors.New("TODO")
	}
	bob.WriteString(fmt.Sprintf(li.PercentFormat, bob2.String()))
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
