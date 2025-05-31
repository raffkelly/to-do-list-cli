package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	CurrentUser UserData            `json:"current_user"`
	Users       map[string]UserData `json:"users"`
}

type UserData struct {
	JWT          string `json:"jwt"`
	RefreshToken string `json:"refresh_token"`
	Username     string `json:"username"`
}

func loadConfig() (Config, error) {
	//initialize config struct
	var config Config
	config.Users = make(map[string]UserData)

	//get filepath for config file
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}
	configDir := filepath.Join(appDataDir, "ToDoApp")
	configFilePath := filepath.Join(configDir, "config.json")

	//check if config file exists
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return config, nil // Return default config
	}

	//read from config file and return information
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)

	//set base URL
	return config, err
}

func saveConfig(config Config) error {
	//get config filepath
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(appDataDir, "ToDoApp")
	configFilePath := filepath.Join(configDir, "config.json")

	//make the config directory if needed
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	//marshal config data and then write to file
	data, err := json.MarshalIndent(config, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, data, 0644)
}
