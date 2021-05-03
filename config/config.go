package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	IDSalt string `json:"salt"`
	Auth   struct {
		JWTPassword string `json:"jwt_password"`
	} `json:"auth"`
	Services struct {
		DB struct {
			URL      string `json:"url"`
			Port     string `json:"port"`
			User     string `json:"user"`
			DBName   string `json:"db_name"`
			Password string `json:"password"`
		} `json:"db"`
		Redis struct {
			URL string `json:"url"`
		} `json:"redis"`
		NLPService struct {
			URL string `json:"url"`
		} `json:"nlp_service"`
	} `json:"services"`
}

var configGlobal *Config

func GetConfig(configPath string) (*Config, error) {
	configFileAbsPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	file, err := os.Open(configFileAbsPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	configFile, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	configGlobal = config
	return config, nil
}

func RetreiveConfig() *Config {
	return configGlobal
}
