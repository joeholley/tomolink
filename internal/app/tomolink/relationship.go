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
	"github.com/sirupsen/logrus"
)

type relationship struct {
	Direction    string   `json:"direction"`
	Relationship string   `json:"relationship"`
	Delta        int      `json:"delta"`
	UUIDs        []string `json:"uuids"`
}

func relationshipLogger(vars *relationship) *logrus.Entry {
	// Populate fields in the logrus structured logging
	logger := hnLog.WithFields(logrus.Fields{
		"relationship": vars.Relationship,
		"delta":        vars.Delta,
		"direction":    vars.Direction,
	})
	if vars.Direction == "unary" ||
		vars.Direction == "discrete" ||
		vars.Direction == "one" ||
		vars.Direction == "single" {
		// Updating this relationship only in source -> target direction, so
		// reflect that in field names
		logger = logger.WithFields(logrus.Fields{
			"UUIDSource": vars.UUIDs[0],
			"UUIDTarget": vars.UUIDs[1],
		})
	} else {
		// All users will have this relationship updated wrt each other, so
		// there is no 'source' or 'target'
		logger = logger.WithFields(logrus.Fields{
			"UUIDs": vars.UUIDs,
		})
	}

	return logger
}
