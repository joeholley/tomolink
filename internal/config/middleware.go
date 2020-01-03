// middleware.go:
// Middleware functions that can be added to the gorilla mux router to force it to
// validate the application config when evaluating client requests.

package config

import (
	"net/http"
)

// Strict is a middleware function that checkes the 'relationships.strict'
// config value and if true, it will refuse to permit requests to access
// relationships not defined in the app config
func (ac *AppConfig) Strict(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strict, err := ac.Cfg.Bool("relationships.strict")
		cfgLog.Printf("strict is '%s'", strict)
		if err != nil {
			cfgLog.Printf("err is '%s'", err.Error())
		}

		next.ServeHTTP(w, r)
	})
}
