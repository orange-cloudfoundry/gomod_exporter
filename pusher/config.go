package main

import (
	"github.com/orange-cloudfoundry/gomod_exporter/common"
)

// Config -
type Config struct {
	common.BaseConfig
	gwURL         string
	gwSkipSSL     bool
	metricJobName string
	metricNS      string
}

// Validate - Validate configuration object
func (c *Config) Validate() error {
	if err := c.BaseConfig.Validate(); err != nil {
		return err
	}
	return nil
}

// NewConfig - Creates and validates config from given reader
func NewConfig() *Config {
	project := common.GitConfig{
		URL: *projectURL,
		Dir: *projectDir,
	}
	if projectUsername != nil && projectPassword != nil {
		project.Auth = &common.GitAuth{
			Username: *projectUsername,
			Password: *projectPassword,
		}
	}

	config := &Config{
		gwURL:         *pushGwURL,
		gwSkipSSL:     *pushGwSkipSSL,
		metricJobName: *metricJobName,
		metricNS:      *metricNS,
	}

	config.BaseConfig = common.BaseConfig{
		Projects: []common.GitConfig{project},
		Log: common.LogConfig{
			JSON:    *logJSON,
			Level:   *logLevel,
			NoColor: *logNoColor,
		},
	}

	return config
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
