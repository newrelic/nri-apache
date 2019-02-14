// +build integration

package integration

import (
	"flag"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/nri-apache/tests/integration/helpers"
	"github.com/newrelic/nri-apache/tests/integration/jsonschema"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var (
	defaultContainer = "integration_nri-apache_1"

	defaultBinPath   = "/nr-apache"
	defaultStatusURL = "http://127.0.0.1/server-status?auto"

	// cli flags
	container = flag.String("container", defaultContainer, "container where the integration is installed")
	binPath   = flag.String("bin", defaultBinPath, "Integration binary path")

	statusURL = flag.String("status_url", defaultStatusURL, "apache status url")
)

func TestMain(m *testing.M) {
	flag.Parse()

	result := m.Run()

	os.Exit(result)
}

// Returns the standard output, or fails testing if the command returned an error
func runIntegration(t *testing.T, envVars ...string) (string, string, error) {
	t.Helper()

	command := make([]string, 0)
	command = append(command, *binPath)

	var found bool

	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "STATUS_URL") {
			found = true
			break
		}
	}

	if !found && statusURL != nil {
		command = append(command, "--status_url", *statusURL)
	}

	stdout, stderr, err := helpers.ExecInContainer(*container, command, envVars...)

	if stderr != "" {
		log.Debug("Integration command Standard Error: ", stderr)
	}

	return stdout, stderr, err
}

func TestApacheIntegration(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr, err := runIntegration(t, fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NoError(t, err, "Unexpected error")
	assert.NotNil(t, stderr, "unexpected stderr")

	schemaPath := filepath.Join("json-schema-files", "apache-schema.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of Apache integration doesn't have expected format.")
}

func TestApacheIntegrationInvalidStatusURL(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr, err := runIntegration(t, "STATUS_URL=invalidurl", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	expectedErrorMessage := "unsupported protocol scheme"

	errMatch, _ := regexp.MatchString(expectedErrorMessage, stderr)

	assert.Error(t, err, "Expected error")
	assert.Truef(t, errMatch, "Expected error message: '%s', got: '%s'", expectedErrorMessage, stderr)

	assert.NotNil(t, stdout, "unexpected stdout")
}

func TestApacheIntegrationOnlyMetrics(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr, err := runIntegration(t, "METRICS=true", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NoError(t, err, "There is an error executing the Apache Integration binary")
	assert.NotNil(t, stderr, "unexpected stderr")

	schemaPath := filepath.Join("json-schema-files", "apache-schema-metrics.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of Apache integration doesn't have expected format.")
}

func TestApacheIntegrationOnlyInventory(t *testing.T) {
	testName := helpers.GetTestName(t)

	stdout, stderr, err := runIntegration(t, "INVENTORY=true", fmt.Sprintf("NRIA_CACHE_PATH=/tmp/%v.json", testName))

	assert.NoError(t, err, "There is an error executing the Apache Integration binary")

	assert.NotNil(t, stderr, "unexpected stderr")

	schemaPath := filepath.Join("json-schema-files", "apache-schema-inventory.json")

	err = jsonschema.Validate(schemaPath, stdout)
	assert.NoError(t, err, "The output of Apache integration doesn't have expected format.")
}
