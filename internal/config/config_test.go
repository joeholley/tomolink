package config

import (
	"errors"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {

	// Test config read
	_, err := Read("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Read config")

}

func TestReadEnvVarOverride(t *testing.T) {

	// Test env var override
	os.Setenv("DEV", "true")

	// Test config read
	cfg, err := Read("test")
	if err != nil {
		log.Fatal(err)
	}

	// Get value that was overridden
	dev, err := cfg.Bool("envVarOverride.dev")
	if err != nil {
		log.Fatal(err)
	}

	// get number of overrides processed
	count, err := cfg.Int("envVarOverride.count")
	if err != nil {
		log.Fatal(err)
	}

	// Test that override was processed
	// This is not a perfect test but will catch several obvious bugs
	msg := "Attempted config param override using env var 'DEV' failed"
	assert.Equal(t, dev, true, msg)
	assert.Greater(t, count, 0, msg)

	log.Info("Read config")

}

func TestReadMissingFile(t *testing.T) {

	// Attempt to read a file called "missing_defaults.yaml"
	msg := "File with the name 'missing_defaults.yaml' exists! " +
		"Delete or move this file so that unit tests can validate correct " +
		"program behavior when trying to load a file that doesn't exist"
	// TODO: re-enable once it's working
	// Couldn't get this working in current version of testify
	_ = msg
	//assert.NoFileExists("missing_defaults.yaml", msg)

	_, err := Read("missing")
	if err != nil {
		var pathError *os.PathError

		if errors.As(err, &pathError) {
			log.Info("Missing config file causing a pathError as expected")
		} else {
			log.Fatal(err)
		}

	}
}
