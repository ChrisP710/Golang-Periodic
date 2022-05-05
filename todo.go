package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"strings"

	// "regexp"
	// "bytes"
	// "time"

	todo "github.com/1set/todotxt"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("expected command(s)")
		os.Exit(1)
	}

	if os.Args[1] == "ls" {
		getTasks()
	} else if os.Args[1] == "add" {
		addTask()
	} else if os.Args[1] == "rm" {
		removeTask()
	} else if os.Args[1] == "do" {
		completeTask()
	} else if os.Args[1] == "projects" {
		getProjects()
	} else if os.Args[1] == "tags" {
		getTags()
	}
}

// Gets all tasks (NEED TO PRIORTIZE)
func getTasks() {

	if tasklist, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(tasklist)
	}
}

func addTask() {
	// Create empty Tasklist
	var tasklist todo.TaskList

	// Populates tasklist from file
	if err := tasklist.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}

	//Concatenating user input into single string and setting it to a variable
	userInput := strings.Join(os.Args[2:], " ")

	//Parsing userInput into a Task Struct
	tk, err := todo.ParseTask(userInput)
	if err != nil {
		log.Fatal(err)
	}

	//Adding task to List
	tasklist.AddTask(tk)

	//Writing new list to todo.txt
	err = tasklist.WriteToPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func removeTask() {
	var tasklist todo.TaskList

	// Populates tasklist from file
	if err := tasklist.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}

	// Concatenating user input into single string and setting it to a variable
	userInput := os.Args[2]

	intInput, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatal(err)
	}

	tasklist.RemoveTaskByID(intInput)

	//Writing new list to todo.txt
	err = tasklist.WriteToPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func completeTask() {
	var tasklist todo.TaskList

	// Populates tasklist from file
	if err := tasklist.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}

	// Concatenating user input into single string and setting it to a variable
	userInput := os.Args[2]

	intInput, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatal(err)
	}

	task, err := tasklist.GetTask(intInput)
	if err != nil {
		log.Fatal(err)
	}

	task.Complete()

	//Writing new list to todo.txt
	err = tasklist.WriteToPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func getProjects() {
	var tasklist todo.TaskList

	// Populates tasklist from file
	if err := tasklist.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}

	taskProjects := make([]string, 0, 100)

	// for i := 0; i < len(tasklist); i++ {
	// 	if tasklist[i].Projects != nil {
	// 		fmt.Println(tasklist[i].Projects)
	// 	}
	// }

	for i := 0; i < len(tasklist); i++ {
		if tasklist[i].Projects != nil {
			taskProjects = append(tasklist[i].Projects, taskProjects...)
		}
	}

	taskProjectsDupFree := make([]string, 0, 100)

	for i := 0; i < len(taskProjects); i++ {
		taskProjectsDupFree = removeDupStr(taskProjects)
	}

	for i := 0; i < len(taskProjectsDupFree); i++ {
		fmt.Println(taskProjectsDupFree[i])
	}

}

func getTags() {
	var tasklist todo.TaskList

	// Populates tasklist from file
	if err := tasklist.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	}

	// for i := 0; i < len(tasklist); i++ {
	// 	if len(tasklist[i].AdditionalTags) > 0 {
	// 		fmt.Println(tasklist[i].AdditionalTags)
	// 	}
	// }

	var finalMap = make(map[string]bool)

	for _, tk := range tasklist {
		for key, value := range tk.AdditionalTags {
			mapEntry := fmt.Sprintf("%s:%s", key, value)
			finalMap[mapEntry] = true
		}
	}

	for key, _ := range finalMap {
		fmt.Println(key)
	}
}

func removeDupStr(tasklistSlice []string) []string {
	allKeys := make(map[string]bool)
	projectList := []string{}
	for _, project := range tasklistSlice {
		if _, value := allKeys[project]; !value {
			allKeys[project] = true
			projectList = append(projectList, project)
		}
	}
	return projectList
}
