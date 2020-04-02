package main

import (
	"fmt"
	"github.com/orange-cloudfoundry/gomod_exporter/common"
	"io"
	"os"
	"time"
)

// ExporterConfig -
type ExporterConfig struct {
	Interval  string `yaml:"interval"`
	Path      string `yaml:"path"`
	Namespace string `yaml:"namespace"`

	intervalDuration time.Duration
}

func (c *ExporterConfig) validate() error {
	if 0 == len(c.Path) {
		c.Path = "/metrics"
	}
	if 0 == len(c.Namespace) {
		c.Path = "gomod"
	}
	if 0 == len(c.Interval) {
		c.Interval = "24h"
	}
	val, err := time.ParseDuration(c.Interval)
	if err != nil {
		return fmt.Errorf("invalid exporter.duration value '%s': %s", c.Interval, err)
	}
	c.intervalDuration = val
	return nil
}

// WebConfig -
type WebConfig struct {
	Listen      string `yaml:"listen"`
	SSLKeyPath  string `yaml:"ssl_key"`
	SSLCertPath string `yaml:"ssh_path"`
}

func (c *WebConfig) validate() error {
	if 0 == len(c.Listen) {
		c.Listen = ":23352"
	}
	if 0 != len(c.SSLKeyPath) {
		if _, err := os.Stat(c.SSLKeyPath); err != nil {
			return fmt.Errorf("invalid configuration web.ssl_key: file not found")
		}
	}
	if 0 != len(c.SSLCertPath) {
		if _, err := os.Stat(c.SSLCertPath); err != nil {
			return fmt.Errorf("invalid configuration web.ssl_cert: file not found")
		}
	}
	return nil
}

// Config -
type Config struct {
	common.BaseConfig
	Exporter ExporterConfig `yaml:"exporter"`
	Web      WebConfig      `yaml:"web"`
}

// Validate - Validate configuration object
func (c *Config) Validate() error {
	if err := c.BaseConfig.Validate(); err != nil {
		return err
	}
	if err := c.Exporter.validate(); err != nil {
		return fmt.Errorf("invalid exporter configuration: %s", err)
	}
	if err := c.Web.validate(); err != nil {
		return fmt.Errorf("invalid exporter configuration: %s", err)
	}
	return nil
}

// NewConfig - Creates and validates config from given reader
func NewConfig(file io.Reader) *Config {
	config := Config{}
	if err := common.LoadConfig(file, &config); err != nil {
		os.Exit(1)
	}
	return &config
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
