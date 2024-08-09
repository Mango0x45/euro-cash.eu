package main

import (
	"cmp"
	"errors"
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
	"/":                 templates.Root(),
	"/about":            templates.About(),
	"/coins":            templates.Coins(),
	"/coins/designs":    templates.CoinsDesigns(),
	"/coins/designs/nl": templates.CoinsDesignsNl(),
	"/language":         templates.Language(),
}

func main() {
	i18n.InitPrinters()

	port := flag.Int("port", 8080, "port number")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("GET /designs/", fs)
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /fonts/", fs)
	mux.Handle("GET /style.css", fs)
	mux.Handle("GET /", middleware.I18n(http.HandlerFunc(finalHandler)))
	mux.Handle("POST /language", http.HandlerFunc(setUserLanguage))

	portStr := ":" + strconv.Itoa(*port)
	log.Println("Listening on", portStr)
	log.Fatal(http.ListenAndServe(portStr, mux))
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	p := r.Context().Value("printer").(i18n.Printer)

	/* Strip trailing slash from the URL */
	path := r.URL.Path
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if c, ok := components[path]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, p.T("Page not found"))
	} else {
		/* When a user clicks on the language button to be taken to the
		   language selection page, we need to set a redirect cookie so
		   that after selecting a language they are taken back to the
		   original page they came from. */
		if path == "/language" {
			http.SetCookie(w, &http.Cookie{
				Name:  "redirect",
				Value: cmp.Or(r.Referer(), "/"),
			})
		}
		templates.Base(nil, c).Render(r.Context(), w)
	}
}

func setUserLanguage(w http.ResponseWriter, r *http.Request) {
	loc := r.FormValue("locale")
	_, ok := i18n.Printers[strings.ToLower(loc)]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Locale ‘%s’ is invalid or unsupported", loc)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "locale",
		Value:  loc,
		MaxAge: math.MaxInt32,
	})

	if c, err := r.Cookie("redirect"); errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:   "redirect",
			MaxAge: -1,
		})
		http.Redirect(w, r, c.Value, http.StatusFound)
	}
}
