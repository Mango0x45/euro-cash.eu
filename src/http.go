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
	"strings"
	"time"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
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
		mux.Handle("GET /style-2.css", fs)
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
			Name:    "redirect",
			Value:   cmp.Or(r.Referer(), "/"),
			Expires: time.Now().Add(24 * time.Hour),
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

		if p == pZero {
			td.Printer = bestFitLanguage(r.Header.Get("Accept-Language"))
			http.SetCookie(w, &http.Cookie{
				Name:  "redirect",
				Value: r.URL.Path,
			})
			templates["/language"].Execute(w, td)
		} else {
			td.Printer = p
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

		td.Type = r.FormValue("type")
		switch td.Type {
		case "circ", "nifc", "proof":
		default:
			td.Type = "circ"
		}

		td.FilterBy = r.FormValue("filter-by")
		switch td.FilterBy {
		case "country", "year":
		default:
			td.FilterBy = "country"
		}

		var err error
		mt := dbx.NewMintageType(td.Type)

		switch td.FilterBy {
		case "country":
			td.Code = r.FormValue("country")
			if !slices.ContainsFunc(td.Countries, func(c country) bool {
				return c.Code == td.Code
			}) {
				td.Code = td.Countries[0].Code
			}
			td.CountryMintages, err = dbx.GetMintagesByCountry(td.Code, mt)
		case "year":
			td.Year, err = strconv.Atoi(r.FormValue("year"))
			if err != nil {
				td.Year = 1999
			}
			td.YearMintages, err = dbx.GetMintagesByYear(td.Year, mt)

			/* NOTE: It’s safe to use MustParse() here, because by this
			   point we know that all BCPs are valid. */
			c := collate.New(language.MustParse(td.Printer.Bcp))
			for i, r := range td.YearMintages.Standard {
				name := td.Printer.GetC(
					countryCodeToName[r.Country], "Place Name")
				td.YearMintages.Standard[i].Country = name
			}
			for i, r := range td.YearMintages.Commemorative {
				name := td.Printer.GetC(
					countryCodeToName[r.Country], "Place Name")
				td.YearMintages.Commemorative[i].Country = name
			}
			slices.SortFunc(td.YearMintages.Standard, func(x, y dbx.MSYearRow) int {
				Δ := c.CompareString(x.Country, y.Country)
				if Δ == 0 {
					Δ = c.CompareString(x.Mintmark, y.Mintmark)
				}
				return Δ
			})
			slices.SortFunc(td.YearMintages.Commemorative, func(x, y dbx.MCYearRow) int {
				Δ := c.CompareString(x.Country, y.Country)
				if Δ == 0 {
					Δ = c.CompareString(x.Mintmark, y.Mintmark)
				}
				return Δ
			})
		}
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

func bestFitLanguage(qry string) i18n.Printer {
	type option struct {
		bcp     string
		quality float64
	}
	var xs []option

	for subqry := range strings.SplitSeq(qry, ",") {
		var o option
		subqry = strings.TrimSpace(subqry)
		parts := strings.Split(subqry, ";")
		o.bcp = strings.ToLower(parts[0])
		if len(parts) == 1 {
			o.quality = 1
		} else {
			n, err := fmt.Sscanf(parts[1], "q=%f", &o.quality)
			if n != 1 || err != nil {
				/* Malformed query string; just give up */
				return i18n.DefaultPrinter
			}
		}
		xs = append(xs, o)
	}

	slices.SortFunc(xs, func(x, y option) int {
		return cmp.Compare(y.quality, x.quality)
	})

	for _, x := range xs {
		if p, ok := i18n.Printers[x.bcp]; ok {
			return p
		}
	}
	return i18n.DefaultPrinter
}
