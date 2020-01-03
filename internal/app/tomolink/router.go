// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tomolink is the main application code that sets up the http router
// and functions that handle the http endpoints.
package tomolink

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/joeholley/tomolink/internal/config"
	"github.com/sirupsen/logrus"
)

// String constants to make our path construction more readable.
const source = "{UUIDSource}"
const target = "{UUIDTarget}"

var (
	// Logrus structured logging setup
	tlLog = logrus.WithFields(logrus.Fields{})
)

// Router instantiates a new gorilla mux router and adds the various routes and
// handlers.  This is exported as it is also instantiated by the tests in
// tomolink_test.go
func Router(ac *config.AppConfig) *mux.Router {
	r := mux.NewRouter()
	users := r.PathPrefix("/users").Subrouter()

	/*
		// Loop through all relationships defined in the config file.
		for i := 0; i < config.MaxRelationships; i++ {
			// Figure out config keys for this relationship
			index := fmt.Sprintf("relationships.%d", i)
			nameKey := index + ".name"
			kindKey := index + ".type"

			// Initialize routes and handlers for this relationship
			relationship, err := ac.Cfg.String(nameKey)

			if err != nil {
				if err.Error() == fmt.Sprintf("Required setting '%s' not set", nameKey) {
					tlLog.Info(fmt.Sprintf("No '%s' configured. Done processing relationships from the config", nameKey))
					break
				}
				log.Fatal(err)
			}

			// 'type' is a keyword, so we use var name 'kind' to hold the type
			kind, err := ac.Cfg.String(kindKey)
			if err != nil {
				if err.Error() == fmt.Sprintf("Required setting '%s' not set", kindKey) {
					tlLog.Info(fmt.Sprintf("No '%s' configured. Done processing relationships from the config", kindKey))
					break
				}
				log.Fatal(err)
			}

			// Set up fields for structured logs
			f := logrus.Fields{"relationshipIndex": i,
				"rName":            relationship,
				"relationshipType": kind,
			}

			// Relationship types that support retreiving a single relationship 'score'
			if kind == "timestamp" || kind == "score" {
				// GET endpoint for one score of this relationship type
				route := "/" + source + "/" + relationship + "/" + target
				name := "retrieve_one_" + relationship
				users.Handle(route, Handler{ac, CreateEndpoint}).
					Methods("GET").
					Name("retrieveone" + name + kind)
				tlLog.WithFields(f).Info(fmt.Sprintf("Added '/users%s' route", route))
			}

			// GET endpoint for all of this relationship type of a given user
			route := "/" + source + "/" + relationship
			name := "retrieve_all_" + relationship
			users.Handle(route, Handler{ac, CreateEndpoint}).
				Methods("GET").
				Name(name)
			tlLog.WithFields(f).Info(fmt.Sprintf("Added '/users%s' route", route))

		}
	*/

	// {relationship:(?:1|2|3|4)}
	tlLog.Info(fmt.Sprintf("Relationships: %s", ac.Relationships))

	tlLog.Info("All configured relationship endpoints created")

	// You can find the handlers in handlers.go
	users.Handle("/"+source, Handler{ac, Retrieve}).
		Methods("GET").
		Name("retrieve_user_all")
	//	r.HandleFunc("/users/{UUID}", CreateEndpoint).
	//		Methods("GET", "PUT").
	//		Name("retrieve")
	r.Handle("/create2", Handler{ac, CreateEndpoint2}).Methods("GET").Headers("Content-Type", "application/json").Name("Test2")

	// Middleware handles strict relationship checking if enabled
	// Look for strict relationships flag in the config; default to true
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		r.Use(ac.Strict)
	}

	return r
}
