// middleware.go:
// Middleware functions that can be added to the gorilla mux router to force it to
// validate the application config when evaluating client requests.

package config

import (
	"net/http"
	"reflect"

	"github.com/joeholley/tomolink/internal/models"
	"github.com/sirupsen/logrus"
)

// StrictMW is a middleware function that checkes the 'relationships.strict'
// config value and if true, it will refuse to permit requests to access
// relationships not defined in the app config
func (ac *AppConfig) StrictMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sLog := cfgLog.WithFields(logrus.Fields{
			"relationships.strict": true,
		})

		// Get the request parameters
		var params *models.Relationship
		params = r.Context().Value("params").(*models.Relationship)

		// Check if this is one of the configured relationships
		sLog = sLog.WithFields(logrus.Fields{
			"relationship": params.Relationship,
		})
		valid := params.RelationshipInArray(keys(ac.Relationships))
		sLog = sLog.WithFields(logrus.Fields{
			"valid": valid,
		})
		sLog.Debug("relationship validity check")

		// If it was valid, continue processing the request
		if valid {
			next.ServeHTTP(w, r)
		} else {
			// Otherwise, log and return HTTP 400
			sLog.Warn("failed strict relationship validity check")
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

func keys(a map[string]string) []string {
	j := reflect.ValueOf(a).MapKeys()
	k := make([]string, len(j))
	for i := 0; i < len(a); i++ {
		k[i] = j[i].String()
	}
	return k
}
