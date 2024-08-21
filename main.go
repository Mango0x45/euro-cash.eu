package main

import (
	"cmp"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"git.thomasvoss.com/euro-cash.eu/lib"
	"git.thomasvoss.com/euro-cash.eu/lib/email"
	"git.thomasvoss.com/euro-cash.eu/lib/mintage"
	"git.thomasvoss.com/euro-cash.eu/template"
	"github.com/a-h/templ"
)

var emailDisabled bool

var (
	notFound   = template.NotFound()
	components = map[string]templ.Component{
		"/":                 template.Root(),
		"/about":            template.About(),
		"/coins":            template.Coins(),
		"/coins/designs":    template.CoinsDesigns(),
		"/coins/designs/nl": template.CoinsDesignsNl(),
		"/coins/mintages":   template.CoinsMintages(),
		"/jargon":           template.Jargon(),
		"/language":         template.Language(),
	}
)

func main() {
	lib.InitPrinters()

	port := flag.Int("port", 8080, "port number")
	flag.BoolVar(&emailDisabled, "no-email", false,
		"disables email support")
	flag.StringVar(&email.Config.Host, "smtp-host", "smtp.migadu.com",
		"SMTP server hostname")
	flag.IntVar(&email.Config.Port, "smtp-port", 465,
		"SMTP server port number")
	flag.StringVar(&email.Config.ToAddr, "email-to", "bugs@euro-cash.eu",
		"address to send error messages to")
	flag.StringVar(&email.Config.FromAddr, "email-from", "noreply@euro-cash.eu",
		"address to send error messages from")
	flag.StringVar(&email.Config.Password, "email-password", "",
		"password to authenticate the email client")
	flag.Parse()

	fs := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("GET /designs/", fs)
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /fonts/", fs)
	mux.Handle("GET /style.css", fs)
	mux.Handle("GET /coins/mintages", i18nHandler(mintageHandler(http.HandlerFunc(finalHandler))))
	mux.Handle("GET /", i18nHandler(http.HandlerFunc(finalHandler)))
	mux.Handle("POST /language", http.HandlerFunc(setUserLanguage))

	portStr := ":" + strconv.Itoa(*port)
	log.Println("Listening on", portStr)
	log.Fatal(http.ListenAndServe(portStr, mux))
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	/* Strip trailing slash from the URL */
	path := r.URL.Path
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	c, ok := components[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		c = notFound
	}

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
	template.Base(c).Render(r.Context(), w)
}

func i18nHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p, pZero lib.Printer

		if c, err := r.Cookie("locale"); err == nil {
			p = lib.Printers[strings.ToLower(c.Value)]
		}

		ctx := context.WithValue(
			r.Context(), "printer", cmp.Or(p, lib.DefaultPrinter))

		if p == pZero {
			http.SetCookie(w, &http.Cookie{
				Name:  "redirect",
				Value: r.URL.Path,
			})
			template.Base(template.Language()).Render(ctx, w)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func mintageHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countries := lib.SortedCountries(
			r.Context().Value("printer").(lib.Printer))

		code := strings.ToLower(r.FormValue("code"))
		if !slices.ContainsFunc(countries, func(c lib.Country) bool {
			return c.Code == code
		}) {
			code = countries[0].Code
		}

		ctype := strings.ToLower(r.FormValue("type"))
		switch ctype {
		case "circ", "nifc", "proof":
		default:
			ctype = "circ"
		}

		path := filepath.Join("data", "mintages", code)
		f, err := os.Open(path)
		if err != nil {
			throwError(http.StatusInternalServerError, err, w, r)
			return
		}
		defer f.Close()

		data, err := mintage.Parse(f, path)
		if err != nil {
			throwError(http.StatusInternalServerError, err, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "code", code)
		ctx = context.WithValue(ctx, "type", ctype)
		ctx = context.WithValue(ctx, "mintages", data)
		ctx = context.WithValue(ctx, "countries", countries)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setUserLanguage(w http.ResponseWriter, r *http.Request) {
	loc := r.FormValue("locale")
	_, ok := lib.Printers[strings.ToLower(loc)]
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
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

func throwError(status int, err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(status)
	if emailDisabled {
		log.Print(err)
	} else {
		go func() {
			if err := email.ServerError(err); err != nil {
				log.Print(err)
			}
		}()
	}
	template.Base(template.Error(status)).Render(r.Context(), w)
}
