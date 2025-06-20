package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func run(config *Config) error {
	currentUser, err := config.handleUserAuthentication()
	if err != nil {
		return err
	}
	config.CurrentUser = currentUser
	config.Users[currentUser.Username] = currentUser
	fmt.Println(currentUser.Username, "successfully logged in.")
	time.Sleep(2 * time.Second)
	clearScreen()
	for {
		choice, err := config.displayTaskMenu()
		if err != nil {
			return err
		}
		err = config.handleTaskChoice(choice)
		if err != nil {
			return err
		}
	}
}

func (config *Config) displayTaskMenu() (string, error) {
	_, err := config.displayTaskList()
	if err != nil {
		return "", err
	}
	fmt.Println("===========================================================")
	fmt.Println("1. Add new task")
	fmt.Println("2. Delete task")
	fmt.Println("3. Exit")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Please Enter 1, 2, or 3")
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "1" || input == "2" || input == "3" {
			return input, nil
		}
	}
}

func (config *Config) handleUserAuthentication() (UserData, error) {
	for {
		choice, err := displayMainMenu()
		if err != nil {
			return UserData{}, err
		}
		switch choice {
		case "1":
			//login to existing account
			user, err := config.promptForUserSelection()
			if err != nil {
				return UserData{}, err
			}
			return user, nil
		case "2":
			//create new user
			user, err := config.createNewUser()
			if err != nil {
				return UserData{}, err
			}
			return user, nil
		case "3":
			return UserData{}, errors.New("user exit")
		}
	}
}

func displayMainMenu() (string, error) {
	fmt.Println("Welcome to ToDo")
	fmt.Println("1. Login to existing account")
	fmt.Println("2. Create New Account")
	fmt.Println("3. Exit")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Please Enter 1, 2, or 3")
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "1" || input == "2" || input == "3" {
			return input, nil
		}
	}
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// Fall back to printing newlines for unsupported systems
		fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	}
}
