package main

import "backend/pkg/config"

type IConfig interface {
	GetPort() uint16
	GetPostgresUrl() string
	GetJwtSigningKey() string
	GetKafkaUrl() string
	GetKafkaTopic() string
}

func LoadConfig(filePath string) (IConfig, error) {
	return config.NewFromFile[*Config](filePath)
}

type Config struct {
	Port          uint16 `yaml:"port"`
	PostgresUrl   string `yaml:"postgres_url"`
	JwtSigningKey string `yaml:"jwt_signing_key" validate:"file"`
	KafkaUrl      string `yaml:"kafka_url"`
	KafkaTopic    string `yaml:"kafka_topic"`
}

func (c *Config) GetPort() uint16 {
	return c.Port
}

func (c *Config) GetPostgresUrl() string {
	return c.PostgresUrl
}

func (c *Config) GetJwtSigningKey() string {
	return c.JwtSigningKey
}

func (c *Config) GetKafkaUrl() string {
	return c.KafkaUrl
}

func (c *Config) GetKafkaTopic() string {
	return c.KafkaTopic
}
