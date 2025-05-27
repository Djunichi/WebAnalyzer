package config

import (
	"WebAnalyzer/internal/migration"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type DBConfig struct {
	URL         string        `mapstructure:"url"`
	MaxIdleTime time.Duration `mapstructure:"maxidletime"`
	MaxLifetime time.Duration `mapstructure:"maxlifetime"`
	MaxOpenConn int           `mapstructure:"maxopenconn"`
	MaxIdleConn int           `mapstructure:"maxidleconn"`
}

// Config struct that helps parse config.yaml
type Config struct {
	ENV       string           `mapstructure:"env"`
	DB        DBConfig         `mapstructure:"db"`
	HTTPPort  string           `mapstructure:"http_port"`
	Migration migration.Config `mapstructure:"migration"`
}

// Load Loads configuration file from specified path.
func Load(path string) (*Config, error) {
	var config Config

	viper.SetConfigFile(path + "/config.yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	viper.AutomaticEnv() // check ENV variables

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
