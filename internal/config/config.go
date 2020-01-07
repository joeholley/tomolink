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

// Package config provides config file loading and overridding using env vars.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	goconfig "github.com/zpatrick/go-config"
)

var (
	// Logrus structured logging setup
	cfgLog = logrus.WithFields(logrus.Fields{
		"component": "internal.config",
	})
)

// MaxRelationships is the maximum number of types of relationships tracked.
// Out-of-the-box, Tomolink supports tracking of up to 10 different kinds of relationships.
// This can be changed and Tomolink recompiled, but note that raising this
// value indiscriminantly may have an impact on startup times!
const MaxRelationships = 10

// AppConfig holds the loaded static config, env overrides, and any runtime
// configuration for the app (shared database connections, etc)
type AppConfig struct {
	DB            interface{}
	Cfg           *goconfig.Config
	Overrides     map[string]string
	Relationships map[string]string
}

// Load the application goconfig into a goconfig.Config object
// Looks for defaults in <appName>_defaults.YAML
func (ac *AppConfig) Load(appName string) error {
	logFields := logrus.Fields{"application": appName}
	cfgLog = cfgLog.WithFields(logFields)

	// Read and load YAML file
	defaultsFileName := "./" + appName + "_defaults.yaml"
	d := goconfig.NewYAMLFile(defaultsFileName)
	defaults := goconfig.NewOnceLoader(d)

	// Set up environment mappings
	mappings, err := getMappings(goconfig.NewOnceLoader(d))
	if err != nil {
		return err
	}

	// Read env vars
	// If you're reading this and looking for more info on env var naming,
	// check the comment above the getMappings function.
	env := goconfig.NewEnvironment(mappings)
	envLoader := goconfig.NewOnceLoader(env)
	cfgLog.Info("Environment variables will override default YAML goconfig values if both exist. " +
		"Due to differences between how environment variables and YAML goconfig keys are named, " +
		"please see the documentation for more details on this feature")

	// get config parameters overridden by environment variables, and set a
	// counter of how many overrides were proccessed. Consumers of this package can
	// then see what has been overridden, useful for debugging, logging, and testing.
	err = ac.populateOverrides(goconfig.NewOnceLoader(env))
	if err != nil {
		return err
	}

	// build goconfig out of all sources, last source wins in conflicts
	sources := []goconfig.Provider{defaults, envLoader}
	ac.Cfg = goconfig.NewConfig(sources)
	if err := ac.Cfg.Load(); err != nil {
		return err
		//cfgLog.Fatal(err)
	}
	l := fmt.Sprintf("Read config defaults from %s with overrides from env vars", defaultsFileName)
	cfgLog.Info(l)

	// Populate the relationships
	err = ac.populateRelationships()
	if err != nil {
		return err
	}

	// If dev flag is set, dump all goconfig values to log at startup
	if dev, err := ac.Cfg.BoolOr("dev", false); err != nil {
		return err
	} else if dev == true {
		cfgLog = cfgLog.WithFields(logrus.Fields{"dev": "true"})
		cfgLog.Info("[Dev] Logging all goconfiguration settings")

		settings, err := ac.Cfg.Settings()
		if err != nil {
			return err
		}

		for key, val := range settings {
			cfgLog.Info(fmt.Sprintf("[Dev]  %s = %s", key, val))
		}
	}

	return nil
}

//getMappings retrieves all goconfig keys from the YAML file, and automatically
//generates equivalent environment variable names. Config keys from the YAML
//file are hierarchical dot-concatinated strings and environment varaible names
//upper-case hierarchical underscore-concatinated strings based on the YAML
//goconfig keys.
//
//Example:
// YAML file goconfig key			env var name
// relationships.strict				RELATIONSHIPS_STRICT
func getMappings(defaults *goconfig.OnceLoader) (map[string]string, error) {
	mappings := make(map[string]string)

	// Load defaults from YAML file
	sources := []goconfig.Provider{defaults}
	cfg := goconfig.NewConfig(sources)
	if err := cfg.Load(); err != nil {
		return nil, err
		//cfgLog.Fatal(err)
	}
	settings, err := cfg.Settings()
	if err != nil {
		return nil, err
		//cfgLog.Fatal(err)
	}

	// Loop through goconfig keys in YAML file and make an equivalent env var name
	for key := range settings {
		var k string

		// Special case: the port to run on is called "http.port" in the config file but should
		// always be overridden by the env var named "PORT" if it exists
		if key == "http.port" {
			// This env var mapping is required by the cloud run container runtime contract:
			// " The following environment variables are automatically added to the running containers:
			//   PORT:	The port your HTTP server should listen on."
			// https://cloud.google.com/run/docs/reference/container-contract#env-vars
			mappings["PORT"] = key
		} else {
			// All other cases, replicate the name from the YAML file to an
			// all-caps, underscore-delimited env var name
			k = strings.ReplaceAll(strings.ToUpper(key), ".", "_")
			mappings[k] = key
		}
		cfgLog.WithFields(logrus.Fields{
			"yamlConfigKey": key,
			"envVarName":    k,
		}).Debug("Env var name to YAML goconfig key equivalence set")
	}

	return mappings, nil
}

//populateOverrides makes a copy of the goconfig containing ONLY
//environment variables that are overriding YAML goconfig values, and prints
//those to the log. This is just for ease of use and debugging.
func (ac *AppConfig) populateOverrides(overrides *goconfig.OnceLoader) error {

	// Read app goconfig env vars
	ac.Overrides = map[string]string{}
	sources := []goconfig.Provider{overrides}
	cfg := goconfig.NewConfig(sources)
	err := errors.New("")
	if err := cfg.Load(); err != nil {
		return err
	}
	ac.Overrides, err = cfg.Settings()
	if err != nil {
		return err
	}

	// Log each app goconfig env var that is set
	for key, value := range ac.Overrides {
		k := strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		cfgLog.WithFields(logrus.Fields{
			"yamlConfigKey": key,
			"envVarName":    k,
			"overrideValue": value,
		}).Info("Override for default YAML value found in environment variable!")
	}

	return nil
}

// populateRelationships() parses the defined relationships out of the config
// file and puts them into the AppConfig.Relationships map for simple reference
func (ac *AppConfig) populateRelationships() error {

	ac.Relationships = map[string]string{}

	// Loop through all possibley defined relationships, looking for config data
	for i := 0; i < MaxRelationships; i++ {
		// Figure out config keys for this relationship
		index := fmt.Sprintf("relationships.definitions.%d", i)
		nameKey := index + ".name"
		kindKey := index + ".type"

		// Initialize routes and handlers for this relationship
		relationship, err := ac.Cfg.String(nameKey)

		if err != nil {
			if err.Error() == fmt.Sprintf("Required setting '%s' not set", nameKey) {
				cfgLog.Info(fmt.Sprintf("No '%s' configured. Done processing relationships from the config", nameKey))
				break
			}
			cfgLog.Error(err)
			return err
		}

		// 'type' is a keyword, so we use var name 'kind' to hold the type
		kind, err := ac.Cfg.String(kindKey)
		if err != nil {
			if err.Error() == fmt.Sprintf("Required setting '%s' not set", kindKey) {
				cfgLog.Info(fmt.Sprintf("No '%s' configured. Done processing relationships from the config", kindKey))
				break
			}
			cfgLog.Error(err)
			return err
		}
		ac.Relationships[relationship] = kind
	}

	return nil
}
