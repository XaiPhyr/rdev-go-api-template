package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Server       ServerConfig `yaml:"server"`
		Database     DBConfig     `yaml:"database"`
		JWTSecretKey string       `yaml:"jwt_secret_key"`
	}

	ServerConfig struct {
		Port string `yaml:"port"`
		Env  string `yaml:"env"`
	}

	DBConfig struct {
		URL          string `yaml:"url"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		MaxIdleConns int    `yaml:"max_idle_conns"`
	}
)

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return nil, err
	}

	config.Server.Port = getEnv("SERVER_PORT", config.Server.Port)
	config.Server.Env = getEnv("APP_ENV", config.Server.Env)
	config.Database.URL = getEnv("DB_URL", config.Database.URL)
	config.JWTSecretKey = getEnv("JWT_SECRET", config.JWTSecretKey)

	return config, err
}

func getEnv(key, currentVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return currentVal
}
