package main

import (
	"fmt"
	"os"

	"github.com/pmescheriakov/Habit-Tracker/internal/interface/cli"
)

// main
//
//	Entity <user> actions:
//	- users: show users, active status and count of active tracking habits
//	- add: add new user
//	- stop: active to inactive status to user
//	- change: change user info
//	- switch: switch current user
//	- info: info about user's nums habits, each streak and other
//
//	Entity <habit> actions:
//	- add: add new habit to user
//	- stop: stop active user's habit
//	- change: change user's habit
//	- done: mark today-done habit to user
//	- undone: unmark today-done habit to user
func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No args! --> Use just <help> entity-arg to see all possible entities!")
	}

	switch args[0] {
	case "user":
		entityArgs := make([]string, len(args[1:]))
		copy(entityArgs, args[1:])
		cli.HandleUser(entityArgs)
	case "habit":
		entityArgs := make([]string, len(args[1:]))
		copy(entityArgs, args[1:])
		cli.HandleHabit(entityArgs)
	case "help":
		cli.HandleHelp()
	default:
		fmt.Println("Wrong entity! --> Use just <help> entity-arg to see all possible entities!")
	}
}
