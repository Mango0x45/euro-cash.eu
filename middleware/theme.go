package middleware

import (
	"cmp"
	"context"
	"math"
	"net/http"
)

const defaultTheme = "dark"

func Theme(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userTheme string

		/* Grab the ‘theme’ cookie to figure out what the users current
		   theme is and add it to the context.  If the user doesn’t yet
		   have a theme cookie or the cookie they have contains an
		   invalid theme then we fallback to the default theme and
		   (re)set the cookie. */

		c, err := r.Cookie("theme")
		if err == nil {
			switch c.Value {
			case "dark", "light":
				userTheme = c.Value
			}
		}

		theme := cmp.Or(userTheme, defaultTheme)
		if userTheme == "" {
			http.SetCookie(w, &http.Cookie{
				Name:   "theme",
				Value:  theme,
				MaxAge: math.MaxInt32,
			})
		}

		ctx := context.WithValue(r.Context(), "theme", theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
