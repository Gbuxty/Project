package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string          `mapstructure:"env"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	Grpc     *GRPCConfig     `mapstructure:"grpc"`
	Auth     *AuthConfig     `mapstructure:"auth"`
}
type PostgresConfig struct {
	StoragePath string `mapstructure:"storage_path"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type AuthConfig struct {
	SecretKey string `mapstructure:"secret_key"`
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
		return nil, fmt.Errorf("Failed read to configurefile viper:%w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("Failed to Unmurshal configurefile viper:%w", err)
	}

	return &cfg, nil
}
