package pgconfig

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  string
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Database, c.Options)
}

func InitFromEnv() (Config, error) {
	var (
		config Config
		err    error
	)

	if config.Host, err = requireEnvVariable("POSTGRES_HOST"); err != nil {
		return Config{}, err
	}
	if config.Port, err = requireEnvVariable("POSTGRES_PORT"); err != nil {
		return Config{}, err
	}
	if config.User, err = requireEnvVariable("POSTGRES_USER"); err != nil {
		return Config{}, err
	}
	if config.Password, err = requireEnvVariable("POSTGRES_PASS"); err != nil {
		return Config{}, err
	}
	if config.Database, err = requireEnvVariable("POSTGRES_DB"); err != nil {
		return Config{}, err
	}

	config.Options = os.Getenv("POSTGRES_OPTIONS")

	return config, nil
}

func requireEnvVariable(name string) (string, error) {
	var (
		value string
		ok    bool
	)
	if value, ok = os.LookupEnv(name); !ok {
		return "", fmt.Errorf("env variable %v not set", name)
	}
	return value, nil
}
