package main

import (
	"fmt"
	"os"
)

const baseURL = "http://localhost:8080/api/"

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	fmt.Println("INITIALLY LOGGED IN USER:", config.CurrentUser.Username)
	err = run(&config)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	saveErr := saveConfig(config)
	if saveErr != nil {
		fmt.Println("Error saving config file:", saveErr)
	}
}
