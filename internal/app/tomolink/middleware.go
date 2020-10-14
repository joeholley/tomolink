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

// middleware.go:
// Middleware functions that can be added to the gorilla mux router to force it to
// validate the application config when evaluating client requests.

package tomolink

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joeholley/tomolink/internal/json"
	"github.com/joeholley/tomolink/internal/models"
	"github.com/sirupsen/logrus"
)

// normalizeRequestParams is a middleware function that looks for the
// input parameters in the client request, and puts them into the request
// context.  This serves two functions:
//  1) the request body can only be read once. By putting these parameters into
//     the request context object, multiple middleware functions and the HTTP
//     handlers can all access these parameters, instead of only the first
//     function that reads the request body.
//  2) This allows us to take input parameters both through the URI and the
//     request body.
func normalizeRequestParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var urlParams, jsonBodyParams models.Relationship

		//  parse the mux.vars into params
		if len(mux.Vars(r)) > 0 {
			urlParams.UUIDSource = mux.Vars(r)["UUIDSource"]
			if mux.Vars(r)["UUIDTarget"] != "" {
				urlParams.UUIDTarget = mux.Vars(r)["UUIDTarget"]
			}
			if mux.Vars(r)["relationship"] != "" {
				urlParams.Relationship = mux.Vars(r)["relationship"]
			}

			tlLog.Debug("parsed request URI into request context")
		}

		// Decode the JSON body
		err := json.DecodeJSONBody(w, r, &jsonBodyParams)
		if err != nil && err.Error() != "Request body is empty" {
			tlLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("unable to parse request JSON body into params")
		}

		// Merge the parameters specified in the JSON body and the URL.
		// In the case that a request defines a value for the same parameter in both
		// the URI /and/ the request body JSON, an error is produced. The client has to
		// choose one or the other; having both would require defining the behaviour
		// and if misunderstood could cause bad behaviour (overwriting/deleting data!)
		params, err := urlParams.Merge(&jsonBodyParams)
		if err != nil {
			tlLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("unable to parse params, do you have conflicting values?")
		}
		ctx := context.WithValue(r.Context(), "params", params)
		tlLog.WithFields(logrus.Fields{
			"url":   urlParams,
			"json":  jsonBodyParams,
			"rel":   params.Relationship,
			"del":   params.Delta,
			"uuids": params.UUIDSource,
			"uuidt": params.UUIDTarget,
		}).Debug("parsed request parameters into request context")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
