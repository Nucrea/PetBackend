package main

import "backend/pkg/config"

type IConfig interface {
	GetHttpPort() uint16
	GetGrpcPort() uint16
	GetPostgresUrl() string
}

func LoadConfig(filePath string) (IConfig, error) {
	return config.NewFromFile[*Config](filePath)
}

type Config struct {
	HttpPort    uint16 `yaml:"http_port" validate:"required"`
	GrpcPort    uint16 `yaml:"grpc_port" validate:"required"`
	PostgresUrl string `yaml:"postgres_url" validate:"required"`
}

func (c *Config) GetHttpPort() uint16 {
	return c.HttpPort
}

func (c *Config) GetGrpcPort() uint16 {
	return c.GrpcPort
}

func (c *Config) GetPostgresUrl() string {
	return c.PostgresUrl
}
