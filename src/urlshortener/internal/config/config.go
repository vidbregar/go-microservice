// Package config includes configuration structs and New function
// for reading provided configFile yaml and secrets yaml located in secretsDir
//
// Precedence order. Each item takes precedence over the item below it:
// 	- env variables
// 	- secrets yaml
// 	- config yaml
// 	- default
//
// env variable example: REDIS_PORT=1234
package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Server Server `mapstructure:"SERVER" yaml:"server"`
	Logger Logger `mapstructure:"LOGGER" yaml:"logger"`
	Redis  Redis  `mapstructure:"REDIS" yaml:"redis"`
}

type Server struct {
	Address string `mapstructure:"ADDRESS" yaml:"address"`
}

type Logger struct {
	Development bool `mapstructure:"DEVELOPMENT" yaml:"development"`
}

type Redis struct {
	Host            string `mapstructure:"HOST" yaml:"host"`
	Port            string `mapstructure:"PORT" yaml:"port"`
	Username        string `mapstructure:"USERNAME" yaml:"username"`
	Password        string `mapstructure:"PASSWORD" yaml:"password"`
	Retries         int    `mapstructure:"RETRIES" yaml:"retries"`
	MinRetryBackoff int    `mapstructure:"MINRETRYBACKOFF" yaml:"minRetryBackoff"`
	MaxRetryBackoff int    `mapstructure:"MAXRETRYBACKOFF" yaml:"maxRetryBackoff"`
}

func defaultConfig() *Config {
	return &Config{}
}

func New(configFile string, secretsDir string) (*Config, error) {
	v := viper.New()

	configBytes, err := yaml.Marshal(defaultConfig())
	if err != nil {
		return nil, err
	}

	// set defaults
	confReader := bytes.NewReader(configBytes)
	v.SetConfigType("yaml")
	if err := v.MergeConfig(confReader); err != nil {
		return nil, err
	}

	// override with config yaml
	v.SetConfigFile(configFile)
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return nil, err
		}
	}

	// add secrets from all secrets yamls
	secrets, err := ioutil.ReadDir(secretsDir)
	if err != nil || len(secrets) == 0 {
		fmt.Printf("No secrets found in %s\n", secretsDir)
	}

	for _, secret := range secrets {
		if !secret.IsDir() {
			secretBytes, err := ioutil.ReadFile(filepath.Join(secretsDir, secret.Name()))
			if err != nil {
				return nil, err
			}

			err = v.MergeConfig(bytes.NewBuffer(secretBytes))
			if err != nil {
				return nil, err
			}
		}
	}

	// override with env variables
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := &Config{}
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
