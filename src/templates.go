package app

import (
	"html/template"
	"io/fs"
	"log"
	"strings"

	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
	"git.thomasvoss.com/euro-cash.eu/pkg/watch"

	"git.thomasvoss.com/euro-cash.eu/src/dbx"
	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

type templateData struct {
	Debugp     bool
	Printer    i18n.Printer
	Printers   map[string]i18n.Printer
	Code, Type string
	Mintages   dbx.MintageData
	Countries  []country
}

var (
	notFoundTmpl *template.Template
	errorTmpl    *template.Template
	templates    map[string]*template.Template
	funcmap      = map[string]any{
		"ifElse":           ifElse,
		"locales":          i18n.Locales,
		"map":              templateMakeMap,
		"safe":             asHTML,
		"toUpper":          strings.ToUpper,
		"tuple":            templateMakeTuple,
		"withTranslation":  withTranslation,
		"withTranslations": withTranslations,
	}
)

func BuildTemplates(dir fs.FS) {
	ents := Try2(fs.ReadDir(dir, "."))
	notFoundTmpl = buildTemplate(dir, "-404")
	errorTmpl = buildTemplate(dir, "-error")
	templates = make(map[string]*template.Template, len(ents))

	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		buildAndSetTemplate(dir, name)
		if Debugp {
			go watch.FileFS(dir, name, func() {
				defer func() {
					if p := recover(); p != nil {
						log.Println(p)
					}
				}()

				buildAndSetTemplate(dir, name)
				log.Printf("Template ‘%s’ updated\n", name)
			})
		}
	}
}

func buildAndSetTemplate(dir fs.FS, name string) {
	path := "/"
	name = strings.TrimSuffix(name, ".html.tmpl")
	switch {
	case name[0] == '-':
		return
	case name != "index":
		path += strings.ReplaceAll(name, "-", "/")
	}
	templates[path] = buildTemplate(dir, name)
}

func buildTemplate(dir fs.FS, name string) *template.Template {
	names := [...]string{"-base", "-navbar", name}
	for i, s := range names {
		names[i] = s + ".html.tmpl"
	}
	return template.Must(template.
		New("-base.html.tmpl").
		Funcs(funcmap).
		ParseFS(dir, names[:]...))
}

func asHTML(s string) template.HTML {
	return template.HTML(s)
}

func templateMakeTuple(args ...any) []any {
	return args
}

func templateMakeMap(args ...any) map[string]any {
	if len(args)&1 != 0 {
		/* TODO: Handle error */
		args = args[:len(args)-1]
	}
	m := make(map[string]any, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		k, ok := args[i].(string)
		if !ok {
			/* TODO: Handle error */
			continue
		}
		m[k] = args[i+1]
	}
	return m
}

func ifElse(b bool, x, y any) any {
	if b {
		return x
	}
	return y
}

func withTranslation(tag, bcp, text string, trans template.HTML,
	spanAttrs ...string) template.HTML {
	name, _, _ := strings.Cut(tag, " ")

	var bob strings.Builder
	bob.WriteByte('<')
	bob.WriteString(tag)
	bob.WriteString(`><span lang="`)
	bob.WriteString(bcp)
	bob.WriteString(`">`)
	bob.WriteString(text)
	bob.WriteString("</span>")

	if text != string(trans) {
		bob.WriteString(`<br><span class="translation"`)
		for _, s := range spanAttrs {
			bob.WriteByte(' ')
			bob.WriteString(s)
		}
		bob.WriteByte('>')
		bob.WriteString(string(trans))
		bob.WriteString("</span>")
	}

	bob.WriteString("</")
	bob.WriteString(name)
	bob.WriteByte('>')
	return template.HTML(bob.String())
}

func withTranslations(tag string, text string, translations ...[]any) template.HTML {
	var bob strings.Builder
	bob.WriteByte('<')
	bob.WriteString(tag)
	bob.WriteByte('>')
	bob.WriteString(text)

	/* TODO: Assert that the pairs are [2]string */
	for _, pair := range translations {
		if text == pair[1] {
			continue
		}
		bob.WriteString(`<br><span class="translation"`)
		if pair[0].(string) != "" {
			bob.WriteString(` lang="`)
			bob.WriteString(pair[0].(string))
			bob.WriteString(`">`)
		} else {
			bob.WriteByte('>')
		}
		bob.WriteString(pair[1].(string))
		bob.WriteString("</span>")
	}

	bob.WriteString("</")
	bob.WriteString(tag)
	bob.WriteByte('>')
	return template.HTML(bob.String())
}

func (td templateData) Get(fmt string, args ...map[string]any) template.HTML {
	return template.HTML(td.Printer.Get(fmt, args...))
}

func (td templateData) GetC(fmt, ctx string, args ...map[string]any) template.HTML {
	return template.HTML(td.Printer.GetC(fmt, ctx, args...))
}

func (td templateData) GetN(fmtS, fmtP string, n int, args ...map[string]any) template.HTML {
	return template.HTML(td.Printer.GetN(fmtS, fmtP, n, args...))
}
