package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CacheControl provides a HTTP middleware to emits Cache-Control headers for
// selected URL prefixes.
func CacheControl(next http.Handler, age int, prefixes ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cacheable := false
		for _, prefix := range prefixes {
			if strings.HasPrefix(r.URL.Path, prefix) {
				cacheable = true
				break
			}
		}

		if cacheable {
			w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(age))
		}

		next.ServeHTTP(w, r)
	})
}
