package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Task    string     `json:"task"`
	DueDate *time.Time `json:"due_date"`
	ID      uuid.UUID  `json:"ID"`
}

func (config *Config) createTask(task string, dueDate *time.Time) error {
	newTask := Task{
		Task:    task,
		DueDate: dueDate,
	}

	jsonTask, err := json.Marshal(newTask)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"tasks", bytes.NewReader(jsonTask))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	var confirmation struct{}
	err = config.doAuthenticatedRequest(http.DefaultClient, req, &confirmation)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) getTasks() ([]Task, error) {
	var taskList []Task
	req, err := http.NewRequest("GET", baseURL+"tasks", nil)
	if err != nil {
		return []Task{}, err
	}
	err = config.doAuthenticatedRequest(http.DefaultClient, req, &taskList)
	if err != nil {
		return []Task{}, err
	}
	return taskList, nil
}

func sortTasks(tasks []Task) {
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].DueDate == nil && tasks[j].DueDate != nil {
			return true
		}
		if tasks[i].DueDate != nil && tasks[j].DueDate == nil {
			return false
		}
		if tasks[i].DueDate != nil && tasks[j].DueDate != nil {
			return tasks[i].DueDate.Before(*tasks[j].DueDate)
		}
		return false
	})
}

func (config *Config) handleTaskChoice(choice string) error {
	switch choice {
	case "1":
		err := config.handleAddTask()
		if err != nil {
			return err
		}
		return nil
	case "2":
		err := config.handleDeleteTask()
		if err != nil {
			return err
		}
	case "3":
		return errors.New("user exit")
	}
	return nil
}

func (config *Config) handleAddTask() error {
	clearScreen()
	var taskText string
	var date *time.Time
	for {
		fmt.Println("Enter task description:")
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if input != "" {
			taskText = input
			break
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter task due date (YYYY-MM-DD) or leave blank if no due date (ctrl+c to exit):")
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		if input == "" {
			date = nil
			break
		}
		dateValue, err := time.Parse("2006-01-02", input)
		if err == nil {
			date = &dateValue
			break
		}
		fmt.Println("Please leave input blank or enter date in the format YYYY-MM-DD.")
	}
	err := config.createTask(taskText, date)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) handleDeleteTask() error {
	clearScreen()
	userTasks, err := config.displayTaskList()
	if err != nil {
		return err
	}
	fmt.Println("========================================================")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter the number of the task you would like to delete (ctrl+c to exit):")
		fmt.Print(">")
		scanner.Scan()
		input := scanner.Text()
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(userTasks) {
			fmt.Println("Invalid Choice")
		} else {
			taskID := userTasks[choice-1].ID.String()
			return config.deleteTask(taskID)
		}

	}
}

func (config *Config) deleteTask(taskID string) error {
	req, err := http.NewRequest("DELETE", baseURL+"tasks/"+taskID, nil)
	if err != nil {
		return err
	}

	//this should be struct with message field i think, not a string
	var confirmation struct {
		Message string `json:"message"`
	}
	err = config.doAuthenticatedRequest(http.DefaultClient, req, &confirmation)
	if err != nil {
		return err
	}
	return nil
}

func (config *Config) displayTaskList() ([]Task, error) {
	//function to display a numbered list of all tasks for a user, followed by a list of options (add task, delete task, change task, exit)
	//return users choice
	userTasks, err := config.getTasks()
	if err != nil {
		return nil, err
	}

	//sort userTasks by due date, starting with null dates
	sortTasks(userTasks)
	if len(userTasks) != 0 {
		for i, task := range userTasks {
			fmt.Printf("%v: %v\n", (i + 1), task.Task)
			if task.DueDate != nil {
				fmt.Printf("	Due on: %v\n", task.DueDate.Format("2006-01-02"))
			}
		}
	}
	return userTasks, nil
}
