package configuration_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/raafvargas/wrapit/configuration"
	"github.com/stretchr/testify/assert"
)

func TestFromYAML(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "configuration.yaml")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	cfg := new(configuration.Config)
	err = configuration.FromYAML(file.Name(), cfg)
	assert.NoError(t, err)
}

func TestFromYAMLWithNoFile(t *testing.T) {
	cfg := new(configuration.Config)
	err := configuration.FromYAML("/nofile", cfg)

	assert.Error(t, err)
}
