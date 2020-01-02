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
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/sirupsen/logrus"
)

var (
	hnLog = logrus.WithFields(logrus.Fields{})
)

// CreateEndpoint2 ++
func CreateEndpoint2(common *Common, w http.ResponseWriter, r *http.Request) error {
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
func CreateEndpoint(common *Common, w http.ResponseWriter, r *http.Request) error {
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
}
