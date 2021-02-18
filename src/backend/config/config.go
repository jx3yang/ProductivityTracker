package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env string `env:"ENV" yaml:"env"`

	Port string `env:"PORT" yaml:"port"`

	DBHost     string `env:"DB_HOST" yaml:"db_host"`
	DBPort     string `env:"DB_PORT" yaml:"db_port"`
	DBUsername string `env:"DB_USERNAME" yaml:"db_username"`
	DBPassword string `env:"DB_PASSWORD" yaml:"db_password"`
}

const Prod = "PROD"

func GetConfig(configFileYml *string) (*Config, error) {
	if configFileYml != nil {
		yamlFile, err := ioutil.ReadFile(*configFileYml)
		if err != nil {
			return nil, err
		}
		var config Config
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			return nil, err
		}
		return &config, nil
	}
	return &Config{
		Env:        os.Getenv("ENV"),
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}, nil
}
