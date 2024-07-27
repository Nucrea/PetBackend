package config

type IConfig interface {
	GetPort() uint16
	GetPostgresUrl() string
	GetJwtSigningKey() string
}

type Config struct {
	Port          uint16 `yaml:"port"`
	PostgresUrl   string `yaml:"postgres_url"`
	JwtSigningKey string `yaml:"jwt_signing_key" validate:"file"`
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
