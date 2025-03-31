package settings

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func GetDatabaseConfig() (DatabaseConfig, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return DatabaseConfig{}, err
	}
	defer file.Close()

	var config DatabaseConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return DatabaseConfig{}, err
	}
	return config, nil
}
