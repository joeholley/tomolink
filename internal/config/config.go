package config

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	goconfig "github.com/zpatrick/go-config"
)

var (
	// Logrus structured logging setup
	cfgLog *log.Entry
)

// Read the application goconfig into a goconfig.Config object
// Looks for defaults in <appName>_defaults.YAML
func Read(appName string) (*goconfig.Config, error) {
	logFields := log.Fields{
		"application": appName,
		"component":   "config",
	}
	cfgLog = log.WithFields(logFields)

	// Read and load YAML file
	defaultsFileName := "./" + appName + "_defaults.yaml"
	d := goconfig.NewYAMLFile(defaultsFileName)
	defaults := goconfig.NewOnceLoader(d)

	// Set up environment mappings
	mappings, err := getMappings(goconfig.NewOnceLoader(d))
	if err != nil {
		return nil, err
		//cfgLog.Fatal(err)
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
	overrides, err := getOverrides(goconfig.NewOnceLoader(env))
	if err != nil {
		return nil, err
	}
	cp := map[string]string{"envVarOverride.count": strconv.Itoa(len(overrides))}
	for k, v := range overrides {
		cp["envVarOverride."+k] = v
	}
	overrideCounter := goconfig.NewStatic(cp)

	// build goconfig out of all sources, last source wins in conflicts
	sources := []goconfig.Provider{defaults, envLoader, overrideCounter}
	cfg := goconfig.NewConfig(sources)
	if err := cfg.Load(); err != nil {
		return nil, err
		//cfgLog.Fatal(err)
	}
	l := fmt.Sprintf("Read config defaults from %s with overrides from env vars", defaultsFileName)
	cfgLog.Info(l)

	// If dev flag is set, dump all goconfig values to log at startup
	if dev, err := cfg.BoolOr("dev", false); err != nil {
		return nil, err
		//cfgLog.Fatal(err)
	} else if dev == true {
		cfgLog = cfgLog.WithFields(log.Fields{"dev": "true"})
		cfgLog.Info("[Dev] Logging all goconfiguration settings")

		settings, err := cfg.Settings()
		if err != nil {
			return cfg, err
			//cfgLog.Fatal(err)
		}

		for key, val := range settings {
			cfgLog.Info(fmt.Sprintf("[Dev]  %s = %s", key, val))
		}
	}

	return cfg, nil
}

//getMappings retrieves all goconfig keys from the YAML file, and automatically
//generates equivalent environment variable names. Config keys from the YAML
//file are hierarchical dot-concatinated strings and environment varaible names
//upper-case hierarchical underscore-concatinated strings based on the YAML
//goconfig keys.
//
//Example:
// YAML file goconfig key			env var name
// tl.relationships.strict		TL_RELATIONSHIPS_STRICT
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
		k := strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		mappings[k] = key
		cfgLog.WithFields(log.Fields{
			"yamlConfigKey": key,
			"envVarName":    k,
		}).Debug("Env var name to YAML goconfig key equivalence set")
	}

	return mappings, nil
}

//logOverrides makes a throw-away copy of the goconfig containing ONLY
//environment variables that are overriding YAML goconfig values, and prints
//those to the log. This is just for ease of use and debugging.
func getOverrides(defaults *goconfig.OnceLoader) (map[string]string, error) {

	// Read app goconfig env vars
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

	// Log each app goconfig env var that is set
	for key, value := range settings {
		k := strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		cfgLog.WithFields(log.Fields{
			"yamlConfigKey": key,
			"envVarName":    k,
			"overrideValue": value,
		}).Info("Override for default YAML value found in environment variable!")
	}

	return settings, nil
}
