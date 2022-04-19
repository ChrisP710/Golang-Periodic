package main

import (
	"fmt"
	"log"
	"os"

	// "flag"
	todo "github.com/1set/todotxt"
)

func main() {

	// 'todo List' subcommand
	// listCmd := flag.NewFlagSet("ls", flag.ExitOnError)

	// inputs for 'todo List' command
	// listAll := listCmd.String("", "", "List All Tasks")
	// listContext := listCmd.String("@context", "", "List Tasks that contain context tag")

	//Ensuring user passed sub command
	if len(os.Args) < 2 {
		fmt.Println("expected subcommand(s)")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "ls": // if its the 'get' command
		getList()
	default:
		{
			fmt.Println("Invalid Input")
		}
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



// func HandleLS(listCmd *flag.FlagSet,  *String,){

// 	listCmd.Parse(os.Args[2:])

// 	if *all == false && *id == "" {
//     	fmt.Print("id is required or specify --all for all videos")
//     	listCmd.PrintDefaults()
// 		os.Exit(1)
// 	}
// }
