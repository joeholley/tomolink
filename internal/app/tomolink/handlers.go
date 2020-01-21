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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/joeholley/tomolink/internal/config"
	"github.com/sirupsen/logrus"
)

var (
	hnLog = logrus.WithFields(logrus.Fields{})
)

// RetrieveUserRelationships handles pulling all information about a single
// user's outgoing relationships from the database and returning it to the HTTP
// client.
func RetrieveUserRelationships(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	reLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return err
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		reLog = params.VerboseLogger()
	}
	reLog.Debug("request parameters retrieved")

	// Get this relationship type for this user
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	docsnap, err := doc.Get(r.Context())
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}

	// Send the results back to the client
	w.Header().Set("Content-Type", "application/json")
	t, err := json.Marshal(docsnap.Data())
	io.WriteString(w, string(t))

	return err

}

// RetrieveSingleRelationship handles pulling information from the database and returning it to the HTTP client.
func RetrieveSingleRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {

	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	reLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		reLog = params.VerboseLogger()
	}
	reLog.Debug("request parameters retrieved")

	// Get this relationship type for this user
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	docsnap, err := doc.Get(r.Context())
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}
	dataMap, err := docsnap.DataAt(params.Relationship)
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}

	// Send the results back to the client
	w.Header().Set("Content-Type", "application/json")
	t, err := json.Marshal(dataMap.(map[string]interface{})[params.UUIDTarget].(int64))
	io.WriteString(w, string(t))

	return err

}

// RetrieveUserRelationshipsByType handles pulling information from the database and returning it to the HTTP client.
func RetrieveUserRelationshipsByType(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {

	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	reLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		reLog = params.VerboseLogger()
	}
	reLog.Debug("request parameters retrieved")

	// Get this relationship type for this user
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	docsnap, err := doc.Get(r.Context())
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}
	dataMap, err := docsnap.DataAt(params.Relationship)
	if err != nil {
		reLog.WithFields(logrus.Fields{"error": err.Error()}).Error("Cannot process client input")
		return fmt.Errorf("Cannot process client input: %w", err.Error())
	}

	// Send the results back to the client
	w.Header().Set("Content-Type", "application/json")
	t, err := json.Marshal(dataMap)
	io.WriteString(w, string(t))

	return err
}

// CreateRelationship ...
func CreateRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	crLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		return err
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		crLog = params.VerboseLogger()
	}
	crLog.Debug("request parameters retrieved")

	// Create the relationship
	newdata := make(map[string]map[string]interface{})
	newdata[params.Relationship] = make(map[string]interface{})
	newdata[params.Relationship][params.UUIDTarget] = params.Delta
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	if params.IsMultipleDirection() {
		// Multiple relationships to create; make a batch
		crLog.Debug("attempting bi-directional relationship batch")
		batch := ac.DB.(*firestore.Client).Batch()
		// First relationship is already ready to go because it's the same as a single directional relationship update
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Add a second write query for the reciprocal relationship
		newdata[params.Relationship] = make(map[string]interface{})
		newdata[params.Relationship][params.UUIDSource] = params.Delta
		doc = ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDTarget)
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Send in the batch update
		_, err := batch.Commit(r.Context())
		if err != nil {
			crLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failure when attempting bi-directional relationship create")
			return err
		}
		crLog.Info("bi-directional relationship created")

	} else {
		// single relationship create
		crLog.Debug("attempting uni-directional relationship create")
		doc.Set(r.Context(),
			newdata,
			firestore.MergeAll)
		crLog.Info("uni-directional relationship created")
	}

	return err
}

// DeleteRelationship ...
func DeleteRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {

	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	drLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		return err
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		drLog = params.VerboseLogger()
	}
	drLog.Debug("request parameters retrieved")

	// Delete the relationship
	newdata := make(map[string]map[string]interface{})
	newdata[params.Relationship] = make(map[string]interface{})
	newdata[params.Relationship][params.UUIDTarget] = firestore.Delete
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	if params.IsMultipleDirection() {
		// Multiple relationships to create; make a batch
		drLog.Debug("attempting bi-directional relationship batch")
		batch := ac.DB.(*firestore.Client).Batch()
		// First relationship is already ready to go because it's the same as a single directional relationship update
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Add a second write query for the reciprocal relationship
		newdata[params.Relationship] = make(map[string]interface{})
		newdata[params.Relationship][params.UUIDSource] = firestore.Delete
		doc = ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDTarget)
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Send in the batch update
		_, err := batch.Commit(r.Context())
		if err != nil {
			drLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failure when attempting bi-directional relationship delete")
			return err
		}
		drLog.Info("bi-directional relationship deleted")

	} else {
		// single relationship delete
		drLog.Debug("attempting uni-directional relationship delete")
		doc.Set(r.Context(),
			newdata,
			firestore.MergeAll)
		drLog.Info("uni-directional relationship deleted")
	}

	return err
}

//UpdateRelationship ...
func UpdateRelationship(ac *config.AppConfig, w http.ResponseWriter, r *http.Request) error {
	// Retrieve request input parameters from Context & validate them
	// This is populated by middleware.go:NormalizeRequestParams()
	urLog := hnLog
	params, err := retrieveAndValidateParameters(ac, r)
	if err != nil {
		return err
	}
	if verbose, _ := ac.Cfg.BoolOr("logging.verbose", true); verbose == true {
		urLog = params.VerboseLogger()
	}
	urLog.Debug("request parameters retrieved")

	// Delete the relationship
	newdata := make(map[string]map[string]interface{})
	newdata[params.Relationship] = make(map[string]interface{})
	newdata[params.Relationship][params.UUIDTarget] = firestore.Increment(params.Delta)
	doc := ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDSource)
	if params.IsMultipleDirection() {
		// Multiple relationships to create; make a batch
		urLog.Debug("attempting bi-directional relationship batch")
		batch := ac.DB.(*firestore.Client).Batch()
		// First relationship is already ready to go because it's the same as a single directional relationship update
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Add a second write query for the reciprocal relationship
		newdata[params.Relationship] = make(map[string]interface{})
		newdata[params.Relationship][params.UUIDSource] = firestore.Increment(params.Delta)
		doc = ac.DB.(*firestore.Client).Collection("users").Doc(params.UUIDTarget)
		batch.Set(doc,
			newdata,
			firestore.MergeAll)

		// Send in the batch update
		_, err := batch.Commit(r.Context())
		if err != nil {
			urLog.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("failure when attempting bi-directional relationship delete")
			return err
		}
		urLog.Info("bi-directional relationship deleted")

	} else {
		// single relationship delete
		urLog.Debug("attempting uni-directional relationship delete")
		doc.Set(r.Context(),
			newdata,
			firestore.MergeAll)
		urLog.Info("uni-directional relationship deleted")
	}

	return err
}
