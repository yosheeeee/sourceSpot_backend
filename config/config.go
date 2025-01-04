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
		Password string `yaml:"password"`
	} `yaml:"database"`
	JWTSecretKey string `yaml:"jwtSecretKey"`
	Port         int    `yaml:"port"`
}

var AppConfig Config

func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Error while opening config file")
	}

	defer file.Close()

	var decoder = yaml.NewDecoder(file)
	if err = decoder.Decode(&AppConfig); err != nil {
		return fmt.Errorf("Error while decoding config file")
	}

	return nil
}
