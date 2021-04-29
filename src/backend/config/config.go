package config

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env string `env:"ENV" yaml:"env"`

	Port string `env:"PORT" yaml:"port"`

	DBName      string   `env:"DB_NAME" yaml:"db_name"`
	DBAddresses []string `env:"DB_ADDRESSES" yaml:"db_addresses"`
	ReplicaSet  string   `env:"REPL_SET" yaml:"repl_set"`
	DBUsername  string   `env:"DB_USERNAME" yaml:"db_username"`
	DBPassword  string   `env:"DB_PASSWORD" yaml:"db_password"`
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
		Env:         os.Getenv("ENV"),
		Port:        os.Getenv("PORT"),
		DBAddresses: strings.Split(os.Getenv("DB_ADDRESSES"), ","),
		ReplicaSet:  os.Getenv("REPL_SET"),
		DBUsername:  os.Getenv("DB_USERNAME"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
	}, nil
}
