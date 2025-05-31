package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Username      string    `json:"username"`
	Token         string    `json:"token"`
	Refresh_Token string    `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (config *Config) promptForUserSelection() (UserData, error) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Username (or ctrl+C to quit):")
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			user, err := config.handleUserLogin(input)
			if err != nil {
				return user, err
			}
			return user, nil
		}
	}
}

func (config *Config) handleUserLogin(username string) (UserData, error) {
	var user UserData
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	//check if username exists in list of previously logged in users
	_, exists := config.Users[username]

	//if it doesn't, login with user's credentials, if it does, attempt to refresh JWT with refresh token
	if !exists {
		clearScreen()
		fmt.Printf("Enter Password for user %v: ", username)
		scanner.Scan()
		password := scanner.Text()
		user, err = config.loginWithCredentials(username, password)
		if err != nil {
			return user, err
		}
	} else {
		user, err = config.refreshUserLogin(username)
		if err != nil {
			return user, err
		}
	}
	//return the successfully logged in userdata struct
	return user, nil
}

func (config *Config) createNewUser() (UserData, error) {
	clearScreen()
	scanner := bufio.NewScanner(os.Stdin)
	var username, password string
	for {
		clearScreen()
		fmt.Print("Enter desired username: ")
		scanner.Scan()
		username = scanner.Text()
		if username == "" {
			fmt.Println("Username can't be blank.")
			continue
		}
		break
	}
	for {
		clearScreen()
		fmt.Print("Create password: ")
		scanner.Scan()
		password = scanner.Text()
		if password == "" {
			fmt.Println("Password can't be blank.")
			continue
		}
		fmt.Print("Retype password: ")
		scanner.Scan()
		checkPassword := scanner.Text()
		if password != checkPassword {
			fmt.Println("Passwords don't match.")
			continue
		}
		break
	}
	newUser := UserDTO{
		Username: username,
		Password: password,
	}
	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return UserData{}, err
	}
	resp, err := http.Post(baseURL+"users", "application/json", bytes.NewReader(userJSON))
	if err != nil {
		return UserData{}, err
	}
	defer resp.Body.Close()

	var createUserResponseStruct User

	err = decodeResponse(resp, &createUserResponseStruct)

	if err != nil {
		return UserData{}, err
	}

	currentUser, err := config.loginWithCredentials(username, password)
	if err != nil {
		return UserData{}, err
	}

	return currentUser, nil
}

func (config *Config) loginWithCredentials(username, password string) (UserData, error) {
	userInfo := UserDTO{
		Username: username,
		Password: password,
	}
	userJSON, err := json.Marshal(userInfo)
	if err != nil {
		return UserData{}, err
	}

	resp, err := http.Post(baseURL+"users/login", "application/json", bytes.NewReader(userJSON))
	if err != nil {
		return UserData{}, err
	}
	defer resp.Body.Close()

	var loginResponseStruct User
	err = decodeResponse(resp, &loginResponseStruct)
	if err != nil {
		return UserData{}, err
	}

	currentUser := UserData{
		JWT:          loginResponseStruct.Token,
		RefreshToken: loginResponseStruct.Refresh_Token,
		Username:     loginResponseStruct.Username,
	}

	return currentUser, nil
}

func (config *Config) refreshUserLogin(username string) (UserData, error) {
	//create http request
	req, err := http.NewRequest("POST", baseURL+"refresh", nil)
	if err != nil {
		return UserData{}, err
	}
	req.Header.Set("Authorization", "Bearer "+config.CurrentUser.RefreshToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserData{}, err
	}
	defer res.Body.Close()

	var currentUser UserData
	newJWT := make(map[string]string)

	if res.StatusCode == 401 {
		scanner := bufio.NewScanner(os.Stdin)
		clearScreen()
		fmt.Printf("Enter Password for user %v: ", username)
		scanner.Scan()
		password := scanner.Text()
		currentUser, err = config.loginWithCredentials(username, password)
		if err != nil {
			return UserData{}, err
		}
	} else {
		err = decodeResponse(res, &newJWT)
		if err != nil {
			return UserData{}, err
		}
		currentUser.JWT = newJWT["token"]
		currentUser.RefreshToken = config.Users[username].RefreshToken
		currentUser.Username = username
	}

	return currentUser, nil
}
