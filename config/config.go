package config

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
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
	JWTSecretKey      string `yaml:"jwtSecretKey"`
	Port              int    `yaml:"port"`
	GitHubOAuthConfig struct {
		ClientID     string   `yaml:"clientId"`
		ClientSecret string   `yaml:"clientSecret"`
		Scopes       []string `yaml:"scopes"`
		RedirectURL  string   `yaml:"redirectUrl"`
	} `yaml:"githubOAuthConfig"`
}

var AppConfig Config

func GetJWTSecretKey() []byte {
	return []byte(AppConfig.JWTSecretKey)
}

func GetGitHubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     AppConfig.GitHubOAuthConfig.ClientID,
		ClientSecret: AppConfig.GitHubOAuthConfig.ClientSecret,
		Scopes:       AppConfig.GitHubOAuthConfig.Scopes,
		RedirectURL:  AppConfig.GitHubOAuthConfig.RedirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
}

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
