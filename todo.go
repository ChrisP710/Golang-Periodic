package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"strings"

	"regexp"
	// "bytes"
	// "time"

	todo "github.com/1set/todotxt"
)

func main() {
	// Could have opted to use switch statements instead of if chain
	tasklist, err := todo.LoadFromPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		fmt.Println("expected command(s)")
		os.Exit(1)
	}

	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "@") == true {
		getTasksContext(tasklist)
	}

	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "+") == true {
		getTasksProjects(tasklist)
	}

	if os.Args[1] == "ls" && len(os.Args[1:]) < 2 {
		getTasksDefault(tasklist)
		// ls |order
	}

	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sort") == true {
		getTasksOrder(tasklist)
		// ls + (projects)
	}

	if os.Args[1] == "add" {
		addTask(tasklist)
	}

	if os.Args[1] == "rm" {
		removeTask(tasklist)
	}

	if os.Args[1] == "do" {
		completeTask(tasklist)
	}

	if os.Args[1] == "projects" {
		getProjects(tasklist)
	}

	if os.Args[1] == "tags" {
		getTags(tasklist)
	}
}

func getTasksDefault(tasklist todo.TaskList) {

	// if tasklist, err := todo.LoadFromPath("todo.txt"); err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Print(tasklist)
	// }

	// Prints completed Tasks for default LS

	tasklistNonCompleted := tasklist.Filter(todo.FilterNotCompleted)
	if err := tasklist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc); err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasklistNonCompleted)

	// Prints non-completed Tasks for default LS | Note: Unsure if completed task are meant to be displayed with default ls command

	tasklistCompleted := tasklist.Filter(todo.FilterCompleted)
	if err := tasklist.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc); err != nil {
		log.Fatal(err)
	}
	fmt.Println(tasklistCompleted)

}

func addTask(tasklist todo.TaskList) {
	// Create empty Tasklist

	// Populates tasklist from file

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

func removeTask(tasklist todo.TaskList) {

	userInput := os.Args[2]

	intInput, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatal(err)
	}

	tasklist.RemoveTaskByID(intInput)

	err = tasklist.WriteToPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func completeTask(tasklist todo.TaskList) {

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

	err = tasklist.WriteToPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}
}

func getProjects(tasklist todo.TaskList) {

	taskProjects := make([]string, 0, 100)

	// for i := 0; i < len(tasklist); i++ {
	// 	if tasklist[i].Projects != nil {
	// 		fmt.Println(tasklist[i].Projects)
	// 	}
	// }

	for i := 0; i < len(tasklist); i++ {
		if (tasklist)[i].Projects != nil {
			taskProjects = append((tasklist)[i].Projects, taskProjects...)
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

func getTags(tasklist todo.TaskList) {

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

// When using function in windows, must put in string due to "|"
func getTasksOrder(tasklist todo.TaskList) {
	// Rudimentary way of checking user input and sorting tasklist (Possibly Revisit)

	// delete below line if you want to include completed or you could opt copy & paste everything with "tasklist = tasklist.Filter(todo.FilterCompleted) to have completed show on another line."

	tasklist = tasklist.Filter(todo.FilterNotCompleted)

	if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttaskidasc") == true {
		if err := tasklist.Sort(todo.SortTaskIDAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttaskiddesc") == true {
		if err := tasklist.Sort(todo.SortTaskIDDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttodotextasc") == true {
		if err := tasklist.Sort(todo.SortTodoTextAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttodotextdesc") == true {
		if err := tasklist.Sort(todo.SortTodoTextDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortpriorityasc") == true {
		if err := tasklist.Sort(todo.SortPriorityAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprioritydesc") == true {
		if err := tasklist.Sort(todo.SortPriorityDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcreateddateasc") == true {
		if err := tasklist.Sort(todo.SortCreatedDateAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcreateddatedesc") == true {
		if err := tasklist.Sort(todo.SortTodoTextAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcompleteddateasc") == true {
		if err := tasklist.Sort(todo.SortCompletedDateAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcompleteddatedesc") == true {
		if err := tasklist.Sort(todo.SortCompletedDateDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortduedateasc") == true {
		if err := tasklist.Sort(todo.SortDueDateAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortduedatedesc") == true {
		if err := tasklist.Sort(todo.SortDueDateDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcontextasc") == true {
		if err := tasklist.Sort(todo.SortContextAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcontextdesc") == true {
		if err := tasklist.Sort(todo.SortContextDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprojectasc") == true {
		if err := tasklist.Sort(todo.SortProjectAsc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprojectdesc") == true {
		if err := tasklist.Sort(todo.SortProjectDesc); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tasklist)
	}
}

func getTasksContext(tasklist todo.TaskList) {

	// fmt.Println(extractContext(strings.ToLower(strings.Join(os.Args, " "))))

	userInputContext := extractContext(strings.Join(os.Args, " "))

	for _, tk := range tasklist {
		for _, ui := range userInputContext {
			// fmt.Println(ui)
			if strings.Contains(tk.Original, ui) {
				fmt.Println(tk.Original)
			}
		}
	}
}

func extractContext(userContext string) []string {
	re := regexp.MustCompile("\\@[a-zA-Z]+")
	contexts := re.FindAllString(userContext, -1)

	return contexts
}

func getTasksProjects(tasklist todo.TaskList) {

	// fmt.Println(extractContext(strings.ToLower(strings.Join(os.Args, " "))))

	userInputProjects := extractProject(strings.Join(os.Args, " "))

	for _, tk := range tasklist {
		for _, ui := range userInputProjects {
			// fmt.Println(ui)
			if strings.Contains(tk.Original, ui) {
				fmt.Println(tk.Original)
			}
		}
	}
}

func extractProject(userProject string) []string {
	re := regexp.MustCompile("\\+[a-zA-Z]+")
	contexts := re.FindAllString(userProject, -1)

	return contexts
}
