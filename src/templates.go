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
	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

type templateData struct {
	Debugp     bool
	Printer    i18n.Printer
	Code, Type string
	Mintages   dbx.MintageData
	Countries  []country
}

var (
	notFoundTmpl *template.Template
	errorTmpl    *template.Template
	templates    map[string]*template.Template
	funcmap      = map[string]any{
		"locales": i18n.Locales,
		"map":     templateMakeMap,
		"safe":    asHTML,
		"sprintf": fmt.Sprintf,
		"toUpper": strings.ToUpper,
		"tuple":   templateMakeTuple,
	}
)

func BuildTemplates(dir fs.FS) {
	ents := Try2(fs.ReadDir(dir, "."))
	notFoundTmpl = buildTemplate(dir, "-404")
	errorTmpl = buildTemplate(dir, "-error")
	templates = make(map[string]*template.Template, len(ents))

	for _, e := range ents {
		name := e.Name()
		buildAndSetTemplate(dir, name)
		if Debugp {
			go watch.FileFS(dir, name, func() {
				defer func() {
					if p := recover(); p != nil {
						log.Print(p)
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

func (td templateData) Get(fmt string, args ...map[string]any) template.HTML {
	return template.HTML(td.Printer.Get(fmt, args...))
}

func (td templateData) GetN(fmtS, fmtP string, n int, args ...map[string]any) template.HTML {
	return template.HTML(td.Printer.GetN(fmtS, fmtP, n, args...))
}
