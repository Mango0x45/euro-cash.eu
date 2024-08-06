package middleware

import (
	"cmp"
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"git.thomasvoss.com/euro-cash.eu/i18n"
	"git.thomasvoss.com/euro-cash.eu/templates"
)

const PrinterKey = "printer"

func I18n(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p, pZero i18n.Printer

		if c, err := r.Cookie("lang"); errors.Is(err, http.ErrNoCookie) {
			log.Println("Language cookie not set")
		} else {
			var ok bool
			p, ok = i18n.Printers[strings.ToLower(c.Value)]
			if !ok {
				log.Printf("Language ‘%s’ is unsupported\n", c.Value)
			}
		}

		ctx := context.WithValue(
			r.Context(), PrinterKey, cmp.Or(p, i18n.DefaultPrinter))

		if p == pZero {
			http.SetCookie(w, &http.Cookie{
				Name:  "redirect",
				Value: r.URL.Path,
			})
			templates.Root(nil, templates.SetLanguage()).Render(ctx, w)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
