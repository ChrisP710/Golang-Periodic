package main

import (
	"fmt"
	"log"
	"os"

	// "regexp"
	// "strconv"
	// "strings"
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
		getList()
	}
}

// Gets all tasks (NEED TO PRIORTIZE)
func getList() {
	if tasklist, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print(tasklist)
	}
}
