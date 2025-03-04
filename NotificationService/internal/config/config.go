package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string        `yaml:"env"`
	Mailer   MailerConfig  `yaml:"mailer"`
	Kafka    KafkaConfig   `yaml:"kafka"`
}

type MailerConfig struct {
	ApiURL    string `yaml:"api_url"`
	ApiToken  string `yaml:"api_token"`
	FromEmail string `yaml:"from_email"`
}

type KafkaConfig struct {
	Broker  string `yaml:"broker"`
	Topic   string `yaml:"topic"`
	GroupID string `yaml:"group_id"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}