package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	// "regexp"
	// "strconv"

	// "time"

	todo "github.com/1set/todotxt"
)

func main() {
	//"Ensuring" user passed command
	if len(os.Args) < 2 {
		fmt.Println("expected command(s)")
		os.Exit(1)
	}

	if os.Args[1] == "ls" {
		getTasks()
	} else if os.Args[1] == "completed" {
		getTasks()
	} else if os.Args[1] == "add" {
		addTask()
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

// Function to add new Task(s) to "todo.txt"
func addTask() {

	f, err := os.OpenFile("todo.txt", os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	userInput := strings.Join(os.Args[2:], " ")

	newLine := "\n"

	addInput := (userInput + newLine)

	_, err2 := f.WriteString(addInput)

	if err2 != nil {
		log.Fatal(err2)
	}

}
