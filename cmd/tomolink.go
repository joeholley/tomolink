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

// Package app contains the application initialization code for the Tomolink service.
package main

import (
	"net/http"

	"github.com/cloudflare/cfssl/log"
	"github.com/joeholley/tomolink/internal/app/tomolink"
	"github.com/joeholley/tomolink/internal/config"
	"github.com/joeholley/tomolink/internal/logging"
	"github.com/sirupsen/logrus"
)

var (
	// Fields to add to the structured logs
	tlLog = logrus.WithFields(logrus.Fields{
		"app":       "tomolink",
		"component": "app.main",
	})
)

func main() {

	// Read config file, file name is
	//   ../internal/config/<string_argument>_defaults.cfg
	// so in this case it will be called
	//   ../internal/config/tomolink_defaults.cfg
	cfg, err := config.Read("tomolink")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load configuration")
	}

	// set up logrus structured logging
	format, err := cfg.StringOr("logging.format", "text")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load logging format configuration")
	}
	lvl, err := cfg.StringOr("logging.level", "info")
	if err != nil {
		tlLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatalf("Cannot load logging level configuration")
	}
	logging.ConfigureLogging(format, lvl)

	// Instantiate router
	router := tomolink.Router()

	log.Info("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":8080", router))
}
