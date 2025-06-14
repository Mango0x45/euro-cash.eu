package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"strings"

	. "git.thomasvoss.com/euro-cash.eu/pkg/try"
	"git.thomasvoss.com/euro-cash.eu/pkg/watch"

	"git.thomasvoss.com/euro-cash.eu/src/dbx"
)

type templateData struct {
	Printer    Printer
	Code, Type string
	Mintages   dbx.MintageData
	Countries  []country
}

var (
	notFoundTmpl *template.Template
	errorTmpl    *template.Template
	templates    map[string]*template.Template
	funcmap      = map[string]any{
		"denoms":  denoms,
		"locales": locales,
		"safe":    asHTML,
		"sprintf": fmt.Sprintf,
		"toUpper": strings.ToUpper,
		"tuple":   templateMakeTuple,
	}
)

func BuildTemplates(dir fs.FS, debugp bool) {
	ents := Try2(fs.ReadDir(dir, "."))
	notFoundTmpl = buildTemplate(dir, "-404")
	errorTmpl = buildTemplate(dir, "-error")
	templates = make(map[string]*template.Template, len(ents))

	for _, e := range ents {
		name := e.Name()
		buildAndSetTemplate(dir, name)
		if debugp {
			go watch.FileFS(dir, name, func() {
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

func denoms() [8]float64 {
	return [8]float64{
		0.01, 0.02, 0.05, 0.10,
		0.20, 0.50, 1.00, 2.00,
	}
}

func locales() []locale {
	return Locales[:]
}

func templateMakeTuple(args ...any) []any {
	return args
}

func (td templateData) T(fmt string, args ...any) string {
	return td.Printer.T(fmt, args...)
}
