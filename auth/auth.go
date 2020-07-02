package auth

import "github.com/gin-gonic/gin"

// Config ...
type Config struct {
	Mock     MockConfig `yaml:"mock"`
	Tenant   string     `yaml:"tenant"`
	JWKS     string     `yaml:"jwks"`
	Audience []string   `yaml:"audience"`
}

// MockConfig ...
type MockConfig struct {
	Enabled       bool                   `yaml:"enabled"`
	Claims        map[string]interface{} `yaml:"claims"`
	Authenticated bool                   `yaml:"authenticated"`
}

// Handler ...
type Handler interface {
	HTTP() gin.HandlerFunc
}

// NewHandler ...
func NewHandler(config *Config) Handler {
	if config.Mock.Enabled {
		return NewMock(
			config.Mock.Authenticated,
			config.Mock.Claims,
		)
	}

	return NewAuth0Handler(config)
}
