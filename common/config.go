package common

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/yaml.v3"
)

// LogConfig -
type LogConfig struct {
	JSON    bool   `yaml:"json"`
	Level   string `yaml:"level"`
	NoColor bool   `yaml:"no_color"`
}

// GitConfig -
type GitConfig struct {
	URL  string   `yaml:"url"`
	Auth *GitAuth `yaml:"auth"`
	Dir  string
}

// GitAuth -
type GitAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *GitConfig) validate() error {
	return nil
}

// Entry - generate log entry for current object
func (c *GitConfig) Entry() *log.Entry {
	username := "(no-auth)"
	if c.Auth != nil && c.Auth.Username != "" {
		username = c.Auth.Username
	}
	return log.WithFields(log.Fields{
		"url":      c.URL,
		"username": username,
	})
}

// AuthMethod - create transport Auth handler
func (c *GitConfig) AuthMethod() transport.AuthMethod {
	if c.Auth == nil {
		return nil
	}

	return &http.BasicAuth{
		Username: c.Auth.Username,
		Password: c.Auth.Password,
	}
}

// Config - Interface
type Config interface {
	Validate() error
}

// BaseConfig -
type BaseConfig struct {
	Log      LogConfig   `yaml:"log"`
	Projects []GitConfig `yaml:"projects"`
}

// Validate - Validate configuration object
func (c *BaseConfig) Validate() error {
	for _, cProject := range c.Projects {
		if err := cProject.validate(); err != nil {
			return fmt.Errorf("invalid bosh configuration: %s", err)
		}
	}
	return nil
}

// LoadConfig - Creates and validates config from given reader
func LoadConfig(file io.Reader, config Config) error {
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("unable to read configuration file : %s", err)
		return err
	}

	if err = yaml.Unmarshal(content, config); err != nil {
		if err = json.Unmarshal(content, &config); err != nil {
			log.Fatalf("unable to read configuration yaml/json file: %s", err)
			return err
		}
	}
	if err = config.Validate(); err != nil {
		log.Fatalf("invalid configuration, %s", err)
		os.Exit(1)
	}
	return nil
}

// Local Variables:
// ispell-local-dictionary: "american"
// End:
