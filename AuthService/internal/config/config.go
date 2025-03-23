package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string          `mapstructure:"env"`
	Postgres *PostgresConfig `mapstructure:"postgres"`
	Auth     *AuthConfig     `mapstructure:"auth"`
	Grpc     *GRPCConfig     `mapstructure:"grpc"`
	Kafka    *KafkaConfig    `mapstructure:"kafka"`
	Redis    *RedisConfig    `mapstructure:"redis"`
}

type AuthConfig struct {
	SecretKey       string        `mapstructure:"secret_key"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type PostgresConfig struct {
	StoragePath string `mapstructure:"storage_path"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type KafkaConfig struct {
	Broker string `mapstructure:"broker"`
	Topic  string `mapstructure:"topic"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
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
