package configuration

import (
	"io/ioutil"
	"path/filepath"

	"github.com/raafvargas/wrapit/api"
	"github.com/raafvargas/wrapit/mongodb"
	"github.com/raafvargas/wrapit/rabbitmq"
	"github.com/raafvargas/wrapit/tracing"
	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	API      *api.Config            `yaml:"api"`
	Mongo    *mongodb.MongoConfig   `yaml:"mongo"`
	Tracing  *tracing.Config        `yaml:"tracing"`
	RabbitMQ *rabbitmq.RabbitConfig `yaml:"rabbitmq"`
}

// FromYAML ...
func FromYAML(file string, dist interface{}) error {
	filename, _ := filepath.Abs(file)

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, dist)
}
