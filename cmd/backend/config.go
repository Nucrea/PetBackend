package main

import "backend/pkg/config"

type IConfig interface {
	GetPort() uint16
	GetPostgresUrl() string
	GetKafkaUrl() string
	GetKafkaTopic() string
}

func LoadConfig(filePath string) (IConfig, error) {
	return config.NewFromFile[*Config](filePath)
}

type Config struct {
	Port        uint16 `yaml:"port"`
	PostgresUrl string `yaml:"postgres_url"`
	KafkaUrl    string `yaml:"kafka_url"`
	KafkaTopic  string `yaml:"kafka_topic"`
}

func (c *Config) GetPort() uint16 {
	return c.Port
}

func (c *Config) GetPostgresUrl() string {
	return c.PostgresUrl
}

func (c *Config) GetKafkaUrl() string {
	return c.KafkaUrl
}

func (c *Config) GetKafkaTopic() string {
	return c.KafkaTopic
}
