package main

import (
	"fmt"
	"log"
	"os"

	"strconv"

	"regexp"
	"strings"

	// Used in some prior implementations
	// "bytes" // "time"

	todo "github.com/1set/todotxt"

	"github.com/TwiN/go-color"
)

func main() {
	tasklist, err := todo.LoadFromPath("todo.txt")
	if err != nil {
		log.Fatal(err)
	}

	tasklist2 := make(map[string]bool)

	// Could have opted to use switch statements instead of if chain

	if len(os.Args) < 2 {
		fmt.Println("expected command(s)")
		os.Exit(1)
	}

	// ******	Bonus	******
	if os.Args[1] == "man" || os.Args[1] == "help" {
		todoExplain()
	}

	// ******   Required Compoments   ******

	// ls - @Context
	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "@") == true {
		getTasksContext(tasklist, &tasklist2)
	}

	// ls - +Project
	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "+") == true {
		getTasksProjects(tasklist, &tasklist2)
	}

	// ls - Default
	if os.Args[1] == "ls" && len(os.Args[1:]) < 2 {
		getTasksDefault(tasklist)
		return
	}

	// ls - |Order
	if os.Args[1] == "ls" && strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sort") == true {
		tasklist3 := getTasksOrder(tasklist, &tasklist2)

		for _, taskNew := range tasklist3 {
			fmt.Println(taskNew.Original)
		}
		return
	}

	// add <task-id>
	if os.Args[1] == "add" {
		addTask(tasklist)
	}

	// remove <task-id>
	if os.Args[1] == "rm" {
		removeTask(tasklist)
	}

	// do <task-id>
	if os.Args[1] == "do" {
		completeTask(tasklist)
	}

	// projects
	if os.Args[1] == "projects" {
		getProjects(tasklist)
	}

	// tags
	if os.Args[1] == "tags" {
		getTags(tasklist)
	}

	for tk, _ := range tasklist2 {
		fmt.Println(tk)
	}
}

//----------------------- Note to self: Avoid using prints throughout | Only when necessary -----------------------------

func getTasksDefault(tasklist todo.TaskList) {

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

	var taskProjects []string

	// Used to test output when creating | Would normally delete
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

	var taskProjectsDupFree []string

	for i := 0; i < len(taskProjects); i++ {
		taskProjectsDupFree = removeDupStr(taskProjects)
	}

	for i := 0; i < len(taskProjectsDupFree); i++ {
		fmt.Println(taskProjectsDupFree[i])
	}

}

// Function to remove duplicate strings for Project function
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
	// Used to test output when creating | Would normally delete
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

func getTasksOrder(tasklist todo.TaskList, tasklist2 *map[string]bool) todo.TaskList {
	// Rudimentary way of checking user input and sorting tasklist (Possibly Revisit)

	// Used to test output when creating | Would normally delete
	// fmt.Println(len(*tasklist2))
	if len(*tasklist2) > 0 {
		tasklist = todo.NewTaskList()

		for key, _ := range *tasklist2 {
			newTask, err := todo.ParseTask(key)
			if err != nil {
				log.Fatal(err)
			}
			tasklist.AddTask(newTask)
		}
	}

	tasklist = tasklist.Filter(todo.FilterNotCompleted)

	// Could add pipe | symbol to substring in string contains, incase user added a sort type as a parameter when using ls not intended for the |order sub-function
	if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttaskidasc") == true {
		if err := tasklist.Sort(todo.SortTaskIDAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttaskiddesc") == true {
		if err := tasklist.Sort(todo.SortTaskIDDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttodotextasc") == true {
		if err := tasklist.Sort(todo.SortTodoTextAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sorttodotextdesc") == true {
		if err := tasklist.Sort(todo.SortTodoTextDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortpriorityasc") == true {
		if err := tasklist.Sort(todo.SortPriorityAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprioritydesc") == true {
		if err := tasklist.Sort(todo.SortPriorityDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcreateddateasc") == true {
		if err := tasklist.Sort(todo.SortCreatedDateAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcreateddatedesc") == true {
		if err := tasklist.Sort(todo.SortTodoTextAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcompleteddateasc") == true {
		if err := tasklist.Sort(todo.SortCompletedDateAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcompleteddatedesc") == true {
		if err := tasklist.Sort(todo.SortCompletedDateDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortduedateasc") == true {
		if err := tasklist.Sort(todo.SortDueDateAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortduedatedesc") == true {
		if err := tasklist.Sort(todo.SortDueDateDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcontextasc") == true {
		if err := tasklist.Sort(todo.SortContextAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortcontextdesc") == true {
		if err := tasklist.Sort(todo.SortContextDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprojectasc") == true {
		if err := tasklist.Sort(todo.SortProjectAsc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else if strings.Contains(strings.ToLower(strings.Join(os.Args, " ")), "sortprojectdesc") == true {
		if err := tasklist.Sort(todo.SortProjectDesc); err != nil {
			log.Fatal(err)
		}
		return tasklist
	} else {
		return tasklist
	}
}

func getTasksContext(tasklist todo.TaskList, tasklist2 *map[string]bool) {
	// Used to test output when creating | Would normally delete
	// fmt.Println(extractContext(strings.ToLower(strings.Join(os.Args, " "))))

	userInputContext := extractContext(strings.Join(os.Args, " "))

	for _, tk := range tasklist {
		for _, ui := range userInputContext {
			// fmt.Println(ui)
			if strings.Contains(tk.Original, ui) {
				(*tasklist2)[tk.Original] = true
			}
		}
	}
}

func extractContext(userContext string) []string {
	re := regexp.MustCompile("\\@[a-zA-Z]+")
	contexts := re.FindAllString(userContext, -1)

	return contexts
}

func getTasksProjects(tasklist todo.TaskList, tasklist2 *map[string]bool) {
	// Used to test output when creating | Would normally delete
	// fmt.Println(extractContext(strings.ToLower(strings.Join(os.Args, " "))))

	userInputProjects := extractProject(strings.Join(os.Args, " "))

	for _, tk := range tasklist {
		for _, ui := range userInputProjects {
			// fmt.Println(ui)
			if strings.Contains(tk.Original, ui) {
				(*tasklist2)[tk.Original] = true
			}
		}
	}
}

func extractProject(userProject string) []string {
	re := regexp.MustCompile("\\+[a-zA-Z]+")
	contexts := re.FindAllString(userProject, -1)

	return contexts
}

// ******	Bonus	******
func todoExplain() {
	// Added Color just for fun
	println(color.InBold("\n\n****************************************************************************************************************"))
	println(color.InBold("       Command                                                      Command Description"))
	println(color.InCyan("ls <optional-parameters>                     Displays TodoList; Optional Parameters: @context, +project, |order;"))
	println(color.InGreen("add <task-string>                            Adds new task to todo.txt"))
	println(color.InRed("rm <task-id>                                 Remove task from todo.txt by inputted taskID"))
	println(color.InYellow("do <task-id>                                 Mark task completed on todo.txt by inputted taskID"))
	println(color.InPurple("tags                                         Displays all tags - No Duplicates"))
	println(color.InBlue("projects                                     Displays all projects - No Duplicates"))
	println(color.InBold("****************************************************************************************************************\n\n"))
}
