package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/orange-cloudfoundry/gomod_exporter/common"
)

// ExporterConfig -
type ExporterConfig struct {
	Interval  string `yaml:"interval"`
	Path      string `yaml:"path"`
	Namespace string `yaml:"namespace"`

	intervalDuration time.Duration
}

func (c *ExporterConfig) validate() error {
	if len(c.Path) == 0 {
		c.Path = "/metrics"
	}
	if len(c.Namespace) == 0 {
		c.Path = "gomod"
	}
	if len(c.Interval) == 0 {
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
	if len(c.Listen) == 0 {
		c.Listen = ":23352"
	}
	if len(c.SSLKeyPath) != 0 {
		if _, err := os.Stat(c.SSLKeyPath); err != nil {
			return fmt.Errorf("invalid configuration web.ssl_key: file not found")
		}
	}
	if len(c.SSLCertPath) != 0 {
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
