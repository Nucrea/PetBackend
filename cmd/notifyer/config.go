package main

import "backend/pkg/config"

func LoadConfig(filePath string) (Config, error) {
	return config.NewFromFile[Config](filePath)
}

type Config struct {
	App struct {
		LogFile    string `yaml:"logFile"`
		ServiceUrl string `yaml:"serviceUrl"`
	}

	Kafka struct {
		Brokers         []string `yaml:"brokers"`
		Topic           string   `yaml:"topic"`
		ConsumerGroupId string   `yaml:"consumerGroupId"`
	} `yaml:"kafka"`

	SMTP ConfigSMTP `yaml:"smtp"`
}

type ConfigSMTP struct {
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
}
