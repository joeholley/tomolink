// Copyright 2019 Google LLC, with excerpts 2019 Matt Silverlock as noted
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

package tomolink

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joeholley/tomolink/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	hnLog = logrus.WithFields(logrus.Fields{})
)

// RetrieveUser handles pulling information from the database and returning it to the HTTP client.
func RetrieveUserRelationships(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	// Populate fields in the logrus structured logging
	reLog := hnLog.WithFields(logrus.Fields{
		"UUIDSource": vars["UUIDSource"],
	})
	reLog.Debug("RetrieveUserRelationships called")

	w.Header().Set("Content-Type", "application/json")
	return nil
}

// RetrieveUser handles pulling information from the database and returning it to the HTTP client.
func RetrieveSingleRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	// Populate fields in the logrus structured logging
	reLog := hnLog.WithFields(logrus.Fields{
		"UUIDSource": vars["UUIDSource"],
	})
	if vars["UUIDTarget"] != "" {
		reLog = reLog.WithFields(logrus.Fields{
			"UUIDTarget": vars["UUIDTarget"],
		})
	}
	if vars["relationship"] != "" {
		reLog = reLog.WithFields(logrus.Fields{
			"relationship": vars["relationship"],
		})
	}
	reLog.Debug("RetrieveSingleRelationship called")

	w.Header().Set("Content-Type", "application/json")
	return nil
}

// RetrieveUser handles pulling information from the database and returning it to the HTTP client.
func RetrieveUserRelationshipsByType(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	// Populate fields in the logrus structured logging
	reLog := hnLog.WithFields(logrus.Fields{
		"UUIDSource": vars["UUIDSource"],
	})
	if vars["relationship"] != "" {
		reLog = reLog.WithFields(logrus.Fields{
			"relationship": vars["relationship"],
		})
	}

	reLog.Debug("retrieveAllUserRelationships called")

	// Check if strict relationships are enabled, in which case we will only
	// process a relationship request if this relationship is defined in the
	// application config
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
	}

	// No target UUID or relationship provided; retreive the entire database record for this user
	if vars["UUIDTarget"] == "" && vars["relationship"] == "" {

	}

	// Both target UUID and relationship are provided; we are being asked to retreive one relationship

	w.Header().Set("Content-Type", "application/json")
	return nil
}

func DeleteRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {

	// Decode input JSON
	var vars relationship
	err := decodeJSONBody(w, r, &vars)
	if err != nil {
		hnLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Error encountered")
		return err
	}
	drLog := hnLog
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		drLog = relationshipLogger(&vars)
	}
	drLog.Debug("JSON request body decoded successfully!")

	// Delete the relationship

	// Check if strict relationships are enabled, in which case we will only
	// delete if this relationship is defined in the application config
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		drLog.Debug("Strict")
	}

	return nil
}

func CreateRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Decode input JSON
	var vars relationship
	err := decodeJSONBody(w, r, &vars)
	if err != nil {
		hnLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Error encountered")
		return err
	}
	crLog := hnLog
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		crLog = relationshipLogger(&vars)
	}
	crLog.Debug("JSON request body decoded successfully!")

	// Create the relationship

	// Check if strict relationships are enabled, in which case we will only
	// create if this relationship is defined in the application config
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		crLog.Debug("Strict")
	}

	return nil
}

func UpdateRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Decode input JSON
	var vars relationship
	err := decodeJSONBody(w, r, &vars)
	if err != nil {
		hnLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Error encountered")
		return err
	}
	urLog := hnLog
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		urLog = relationshipLogger(&vars)
	}
	urLog.Debug("JSON request body decoded successfully!")

	// Create the relationship

	// Check if strict relationships are enabled, in which case we will only
	// create if this relationship is defined in the application config
	if strict, _ := ac.Cfg.BoolOr("relationships.strict", true); strict == true {
		urLog.Debug("Strict")
	}

	return nil
}

/*

// CreateEndpoint2 ++
func CreateEndpoint2(cfg *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	var p map[string]string
	err := decodeJSONBody(w, r, &p)
	if err != nil {
		hnLog.Fatal(err.Error())
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	t, err := json.Marshal(p)
	io.WriteString(w, string(t))
	hnLog.Info(string(t))

	return err
}

// CreateEndpoint is WIP
//func CreateEndpoint(w http.ResponseWriter, r *http.Request) {
func CreateEndpoint(cfg *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Check that the content header (if set) is application/json
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		hnLog.Info(value)

		if value != "application/json" {
			err := errors.New("Content-Type header is not application/json")
			return StatusError{415, err}
		}
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	var readLimit int64
	readLimit = 500

	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, readLimit))
	if r.Body != nil {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, readLimit))
		//body, err := ioutil.ReadAll(r.Body)
		_ = body
		// https: //stackoverflow.com/questions/32710847/what-is-the-best-way-to-check-for-empty-request-body
		if err != nil {
			hnLog.Printf("Error reading body: %v", err)
			return StatusError{400, err}
		}

		if len(body) > 0 {
			io.WriteString(w, string(body))
			return nil
		}
	}
	io.WriteString(w, `{"body": "nope"}`)
	return nil
}*/
