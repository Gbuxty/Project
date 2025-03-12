package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string        `mapstructure:"env"`
	Mailer *MailerConfig `mapstructure:"mailer"`
	Kafka  *KafkaConfig  `mapstructure:"kafka"`
}

type MailerConfig struct {
	ApiURL    string `mapstructure:"api_url"`
	ApiToken  string `mapstructure:"api_token"`
	FromEmail string `mapstructure:"from_email"`
}

type KafkaConfig struct {
	Broker  string `mapstructure:"broker"`
	Topic   string `mapstructure:"topic"`
	GroupID string `mapstructure:"group_id"`
}
func InitFlags()string{
	var configPath string
	flag.StringVar(&configPath,"c","config/local.yaml","Path to the configuration file")
	flag.Parse()
	return configPath
}

func LoadConfig(path string) (*Config, error) {
    viper.SetConfigFile(path)
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("Failed to read config file: %w", err)
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("Failed to unmarshal config: %w", err)
    }

    return &cfg, nil
}