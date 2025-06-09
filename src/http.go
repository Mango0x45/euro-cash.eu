package src

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
	"strings"

	"git.thomasvoss.com/euro-cash.eu/src/dbx"
	"git.thomasvoss.com/euro-cash.eu/src/email"
)

type middleware = func(http.Handler) http.Handler

func Run(port int) {
	fs := http.FileServer(http.Dir("static"))
	final := http.HandlerFunc(finalHandler)
	mux := http.NewServeMux()

	mwareB := chain(firstHandler, i18nHandler) // [B]asic
	mwareC := chain(mwareB, countryHandler)    // [C]ountry
	mwareM := chain(mwareC, mintageHandler)    // [M]intage

	/* TODO: Put this all in an embed.FS */
	mux.Handle("GET /codes/", fs)
	mux.Handle("GET /designs/", fs)
	mux.Handle("GET /favicon.ico", fs)
	mux.Handle("GET /fonts/", fs)
	mux.Handle("GET /storage/", fs)
	mux.Handle("GET /style.min.css", fs)
	mux.Handle("GET /coins/designs", mwareC(final))
	mux.Handle("GET /coins/mintages", mwareM(final))
	mux.Handle("GET /collecting/crh", mwareC(final))
	mux.Handle("GET /", mwareB(final))
	mux.Handle("POST /language", http.HandlerFunc(setUserLanguage))

	portStr := ":" + strconv.Itoa(port)
	log.Println("Listening on", portStr)
	err := http.ListenAndServe(portStr, mux)
	dbx.DB.Close()
	log.Fatal(err)
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
		ctx := context.WithValue(r.Context(), "td", &templateData{})
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
	t.Execute(w, data)
}

func i18nHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p, pZero Printer

		if c, err := r.Cookie("locale"); err == nil {
			p = printers[strings.ToLower(c.Value)]
		}

		td := r.Context().Value("td").(*templateData)
		td.Printer = cmp.Or(p, defaultPrinter)

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
		td.Code = strings.ToLower(r.FormValue("code"))
		if !slices.ContainsFunc(td.Countries, func(c country) bool {
			return c.Code == td.Code
		}) {
			td.Code = td.Countries[0].Code
		}

		td.Type = strings.ToLower(r.FormValue("type"))
		switch td.Type {
		case "circ", "nifc", "proof":
		default:
			td.Type = "circ"
		}

		var err error
		td.Mintages, err = dbx.GetMintages(td.Code)
		if err != nil {
			throwError(http.StatusInternalServerError, err, w, r)
			return
		}

		processMintages(&td.Mintages, td.Type)
		next.ServeHTTP(w, r)
	})
}

func setUserLanguage(w http.ResponseWriter, r *http.Request) {
	loc := r.FormValue("locale")
	_, ok := printers[strings.ToLower(loc)]
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
	go func() {
		if err := email.ServerError(err); err != nil {
			log.Print(err)
		}
	}()
	errorTmpl.Execute(w, struct {
		Code int
		Msg  string
	}{
		Code: status,
		Msg:  http.StatusText(status),
	})
}

func processMintages(md *dbx.MintageData, typeStr string) {
	var typ int
	switch typeStr {
	case "nifc":
		typ = dbx.TypeNifc
	case "proof":
		typ = dbx.TypeProof
	default:
		typ = dbx.TypeCirc
	}

	md.Standard = slices.DeleteFunc(md.Standard,
		func(x dbx.MSRow) bool { return x.Type != typ })
	md.Commemorative = slices.DeleteFunc(md.Commemorative,
		func(x dbx.MCRow) bool { return x.Type != typ })
	slices.SortFunc(md.Standard, func(x, y dbx.MSRow) int {
		if x.Year != y.Year {
			return x.Year - y.Year
		}
		return strings.Compare(x.Mintmark, y.Mintmark)
	})
	slices.SortFunc(md.Commemorative, func(x, y dbx.MCRow) int {
		if x.Year != y.Year {
			return x.Year - y.Year
		}
		if x.Number != y.Number {
			return x.Number - y.Number
		}
		return strings.Compare(x.Mintmark, y.Mintmark)
	})
}
