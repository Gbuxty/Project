package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env string `mapstructure:"env"`
	AuthServiceAdress string `mapstructure:"auth_service_address"`
	HttpServerAdress string `mapstructure:"http_server_adress"`
}

func InitFlags() string {
	var configPath string
	flag.StringVar(&configPath, "c", "config/local.yaml", "Path to the configuration file")
	flag.Parse()
	return configPath
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &cfg, nil
}
