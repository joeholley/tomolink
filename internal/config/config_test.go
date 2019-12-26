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

	count, err := cfg.Int("configProcessor.envVarOverrideCount")
	if err != nil {
		log.Fatal(err)
	}
	msg := "Attempted config param override using env var 'DEV' failed"
	assert.Greater(t, count, 0, msg)

	log.Info("Read config")

}

func TestReadMissingFile(t *testing.T) {

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
