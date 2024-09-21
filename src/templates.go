package src

import (
	"embed"
	"html/template"
	"strings"

	"git.thomasvoss.com/euro-cash.eu/src/mintage"
)

type templateData struct {
	Printer    Printer
	Code, Type string
	Mintages   mintage.Data
	Countries  []country
}

var (
	//go:embed templates/*.html.tmpl
	templateFS   embed.FS
	notFoundTmpl = buildTemplate("404")
	errorTmpl    = buildTemplate("error")
	templates    = map[string]*template.Template{
		"/":         buildTemplate("index"),
		"/about":    buildTemplate("about"),
		"/jargon":   buildTemplate("jargon"),
		"/language": buildTemplate("language"),
	}
	funcmap = map[string]any{
		"safe":    asHTML,
		"locales": locales,
		"toUpper": strings.ToUpper,
		"tuple":   templateMakeTuple,
	}
)

func buildTemplate(names ...string) *template.Template {
	names = append([]string{"base", "navbar"}, names...)
	for i, s := range names {
		names[i] = "templates/" + s + ".html.tmpl"
	}
	return template.Must(template.
		New("base.html.tmpl").
		Funcs(funcmap).
		ParseFS(templateFS, names...))
}

func asHTML(s string) template.HTML {
	return template.HTML(s)
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
