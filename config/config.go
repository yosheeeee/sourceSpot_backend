package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"string"`
	} `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error while opening config file")
	}

	defer file.Close()

	var cfg Config

	var decoder = yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("Error while decoding config file")
	}

	return &cfg, nil
}
