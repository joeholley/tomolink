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
const relStart = "{relationship"
const relEnd = "}"

var (
	// Logrus structured logging setup
	tlLog    = logrus.WithFields(logrus.Fields{})
	relRegex = ""
)

// Router instantiates a new gorilla mux router and adds the various routes and
// handlers.  This is exported as it is also instantiated by the tests in
// tomolink_test.go
func Router(ac *config.AppConfig) *mux.Router {
	r := mux.NewRouter()
	users := r.PathPrefix("/users").Subrouter()

	// Check if strict relationships are enabled, in which case we will only
	// process a relationship request if this relationship is defined in the
	// application config
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		tlLog.Info("Strict relationships turned ON, parsing config to generate routes")
		// Start a regex unnamed group
		relRegex = ":(?:"
		for rel, _ := range ac.Relationships {
			relRegex = relRegex + rel + "|"
		}
		// Remove last trailing vertical pipe and finish the group
		relRegex = relRegex[:len(relRegex)-1] + ")"

	}
	// This is accomplished by creating a regex that matches using an unnamed
	// group.  In practice, this concatinates three values: relStart + relRegex
	// + relEnd.  relStart and relEnd provide us with a gorilla mux URL
	// variable named "relationship", and relRegex can contain the regex to
	// match only the defined relationship names. For example, if the config sets
	// strict relationships on and then defines four relationships called 'a',
	// 'b', 'c', and 'd', the final concatinated value should look something
	// like this:
	//   {relationship:(?:a|b|c|d)}
	// for more information on the format supported by gorilla mux, see the
	// section starting with "Groups can be used inside patterns, as long
	// as they are non-capturing" on
	// https://godoc.org/github.com/gorilla/mux
	relationship := relStart + relRegex + relEnd

	// Relationship types that support retreiving a single relationship 'score'
	// GET endpoint for one score of this relationship type
	route := "/" + source + "/" + relationship + "/" + target
	users.Handle(route, Handler{ac, RetrieveSingleRelationship}).
		Methods("GET").
		Name("TODO")
	tlLog.WithFields(logrus.Fields{
		"route": fmt.Sprintf("/users%s", route),
	}).Info("Added route")

	// GET endpoint for all of one relationship type of a given user
	route = "/" + source + "/" + relationship
	users.Handle(route, Handler{ac, RetrieveUserRelationshipsByType}).
		Methods("GET").
		Name("TODO2")
	tlLog.WithFields(logrus.Fields{
		"route": fmt.Sprintf("/users%s", route),
	}).Info("Added route")

	// You can find the handlers in handlers.go
	users.Handle("/"+source, Handler{ac, RetrieveUserRelationships}).
		Headers("Content-Type", "application/json").
		Methods("GET").
		Name("TODO3")
	//	r.HandleFunc("/users/{UUID}", CreateEndpoint).
	//		Methods("GET", "PUT").
	//		Name("retrieve")
	//r.Handle("/create2", Handler{ac, CreateEndpoint2}).Methods("GET").Headers("Content-Type", "application/json").Name("Test2")
	tlLog.Info("All configured relationship endpoints created")

	// Middleware handles strict relationship checking if enabled
	// Look for strict relationships flag in the config; default to true
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		r.Use(ac.Strict)
	}

	return r
}
