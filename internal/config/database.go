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
//
// This file contains functions related to configuring the database connection.

package config

import (
	"context"
	"strconv"

	"cloud.google.com/go/firestore"
	//"github.com/joeholley/tomolink/internal/database/memory"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

// Connect to the selected database and make a single client for all requests
// going through this API instance to share.
func (ac *AppConfig) Connect(dbEngine string) error {
	logFields := logrus.Fields{"database.engine": dbEngine}
	dbLog := cfgLog.WithFields(logFields)

	// Read DB settings and put them in the options array
	var options []option.ClientOption
	settings, err := ac.Cfg.Settings()
	if val, ok := settings["database.options.grpc.pool"]; ok {
		optLog := dbLog.WithFields(logrus.Fields{"grpc.pool": val})

		// Convert to int and try to set option
		v, err := strconv.Atoi(val)
		if err != nil {
			// Non-fatal error; just log that we can't set the database client
			// option and continue
			optLog.WithFields(logrus.Fields{
				"error": err.Error()},
			).Warning("Unable to set database option")
		} else {
			options = append(options, option.WithGRPCConnectionPool(v))
			optLog.Info("database option set")
		}
	}

	// Additional database engines could be added as cases in this switch statement
	switch dbEngine {
	case "firestore":
		ac.DB, err = firestore.NewClient(context.Background(),
			settings["database.id"],
			options...)
		if err != nil {
			return err
		}
	case "memory":
		// TODO: NYI, this is for local testing in the future
		//ac.DB = memory.NewClient()

	}

	return nil
}
