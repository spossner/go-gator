package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".config/gator/config.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	if err := write(c); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}

func (c *Config) String() string {
	data, err := json.Marshal(*c)
	if err != nil {
		return ""
	}
	return string(data)
}

func Read() (*Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("error getting config file path: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config json file: %w", err)
	}
	return &config, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error retrieving user home directory: %w", err)
	}
	return dir + "/" + configFileName, nil
}

func write(cfg *Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path: %w", err)
	}
	data, err := json.Marshal(*cfg)
	if err != nil {
		return fmt.Errorf("error creating json: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}
	return nil
}
