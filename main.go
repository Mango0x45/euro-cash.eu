package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"git.thomasvoss.com/euro-cash.eu/i18n"
	"git.thomasvoss.com/euro-cash.eu/middleware"
	"git.thomasvoss.com/euro-cash.eu/templates"
	"github.com/a-h/templ"
)

var components = map[string]templ.Component{
	"/":         templates.Index(),
	"/language": templates.SetLanguage(),
}

func main() {
	i18n.InitPrinters()

	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /style.css", fs)
	mux.Handle("GET /", middleware.I18n(http.HandlerFunc(finalHandler)))
	mux.Handle("POST /language", http.HandlerFunc(setUserLanguage))

	portStr := ":" + strconv.Itoa(*port)
	log.Println("Listening on", portStr)
	log.Fatal(http.ListenAndServe(portStr, mux))
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value(middleware.PrinterKey).(i18n.Printer)

	/* Strip trailing slash from the URL */
	path := r.URL.Path
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if c, ok := components[path]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, p.T("Page not found"))
	} else {
		w.WriteHeader(http.StatusOK)
		templates.Root(nil, c).Render(r.Context(), w)
	}
}

func setUserLanguage(w http.ResponseWriter, r *http.Request) {
	loc := r.FormValue(templates.LocaleKey)
	_, ok := i18n.Printers[strings.ToLower(loc)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Locale ‘%s’ is invalid or unsupported", loc)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "lang",
		Value:  loc,
		MaxAge: math.MaxInt32,
	})
}
