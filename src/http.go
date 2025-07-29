package app

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"slices"
	"strconv"

	. "git.thomasvoss.com/euro-cash.eu/pkg/try"

	"git.thomasvoss.com/euro-cash.eu/src/dbx"
	"git.thomasvoss.com/euro-cash.eu/src/email"
	"git.thomasvoss.com/euro-cash.eu/src/i18n"
)

type middleware = func(http.Handler) http.Handler

func Run(port int) {
	fs := http.FileServer(http.Dir("static"))
	final := http.HandlerFunc(finalHandler)
	mux := http.NewServeMux()

	mwareB := chain(firstHandler, i18nHandler) // [B]asic
	mwareC := chain(mwareB, countryHandler)    // [C]ountry
	mwareM := chain(mwareC, mintageHandler)    // [M]intage

	mux.Handle("GET /codes/", fs)
	mux.Handle("GET /designs/", fs)
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /fonts/", fs)
	mux.Handle("GET /storage/", fs)
	if Debugp {
		mux.Handle("GET /style.css", fs)
	} else {
		mux.Handle("GET /style.min.css", fs)
	}
	mux.Handle("GET /coins/designs", mwareC(final))
	mux.Handle("GET /coins/mintages", mwareM(final))
	mux.Handle("GET /collecting/crh", mwareC(final))
	mux.Handle("GET /", mwareB(final))
	mux.Handle("POST /language", http.HandlerFunc(setUserLanguage))

	portStr := ":" + strconv.Itoa(port)
	log.Println("Listening on", portStr)
	Try(http.ListenAndServe(portStr, mux))
}

func chain(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			next = xs[i](next)
		}
		return next
	}
}

func firstHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "td", &templateData{
			Debugp:   Debugp,
			Printers: i18n.Printers,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func finalHandler(w http.ResponseWriter, r *http.Request) {
	/* Strip trailing slash from the URL */
	path := r.URL.Path
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	t, ok := templates[path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		t = notFoundTmpl
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

	data := r.Context().Value("td").(*templateData)
	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func i18nHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p, pZero i18n.Printer

		if c, err := r.Cookie("locale"); err == nil {
			p = i18n.Printers[c.Value]
		}

		td := r.Context().Value("td").(*templateData)
		td.Printer = cmp.Or(p, i18n.DefaultPrinter)

		if p == pZero {
			http.SetCookie(w, &http.Cookie{
				Name:  "redirect",
				Value: r.URL.Path,
			})
			templates["/language"].Execute(w, td)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func countryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td := r.Context().Value("td").(*templateData)
		td.Countries = sortedCountries(td.Printer)
		next.ServeHTTP(w, r)
	})
}

func mintageHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		td := r.Context().Value("td").(*templateData)
		td.Code = r.FormValue("code")
		if !slices.ContainsFunc(td.Countries, func(c country) bool {
			return c.Code == td.Code
		}) {
			td.Code = td.Countries[0].Code
		}

		td.Type = r.FormValue("type")
		switch td.Type {
		case "circ", "nifc", "proof":
		default:
			td.Type = "circ"
		}

		var err error
		td.Mintages, err = dbx.GetMintages(td.Code, dbx.NewMintageType(td.Type))
		if err != nil {
			throwError(http.StatusInternalServerError, err, w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setUserLanguage(w http.ResponseWriter, r *http.Request) {
	loc := r.FormValue("locale")
	_, ok := i18n.Printers[loc]
	if !ok {
		/* TODO: Make this page pretty? */
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
	go email.Send("Server Error", err.Error())
	errorTmpl.Execute(w, struct {
		Code int
		Msg  string
	}{
		Code: status,
		Msg:  http.StatusText(status),
	})
}
