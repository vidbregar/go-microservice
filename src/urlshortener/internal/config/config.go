// Package config includes configuration structs and New function for reading provided config
//
// Precedence order. Each item takes precedence over the item below it:
// 	- env variables
// 	- config yaml
// 	- default
//
// env variable example: REDIS_PORT=1234
package config

import (
	"bytes"
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
	Host     string `mapstructure:"HOST" json:"host"`
	Port     string `mapstructure:"PORT" json:"port"`
	Retries  int    `mapstructure:"RETRIES" json:"retries"`
	Username string `mapstructure:"USERNAME" json:"username"`
	Password string `mapstructure:"PASSWORD" json:"password"`
}

func defaultConfig() *Config {
	return &Config{}
}

func New(configFile string) (*Config, error) {
	v := viper.New()

	b, err := yaml.Marshal(defaultConfig())
	if err != nil {
		return nil, err
	}

	// set defaults
	confReader := bytes.NewReader(b)
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
