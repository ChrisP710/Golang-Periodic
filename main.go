package main

import (
	"fmt"
	"log"
	// "flag"
	todo "github.com/1set/todotxt"
)

func main() {
	fmt.Println("Hello, world.")
	// ...

	if tasklist, err := todo.LoadFromPath("todo.txt"); err != nil {
		log.Fatal(err)
	} else {
		tasks := tasklist.Filter(todo.FilterNotCompleted).Filter(todo.FilterDueToday, todo.FilterHasPriority)
		_ = tasks.Sort(todo.SortPriorityAsc, todo.SortProjectAsc)
		for i, t := range tasks {
			fmt.Println(t.Todo)
			// oh really?
			tasks[i].Complete()
		}
		if err = tasks.WriteToPath("today-todo.txt"); err != nil {
			log.Fatal(err)
		}
	}
}
