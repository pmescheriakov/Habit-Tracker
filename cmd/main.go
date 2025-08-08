package main

import (
	"fmt"
	"os"

	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
	"github.com/pmescheriakov/Habit-Tracker/internal/interface/cli"
)

// main is the entry point of the Habit Tracker CLI application.
//
//	It parses command-line arguments and delegates handling to the appropriate
//	entity-specific handler: user, habit, or help.
//	If no arguments are provided, it prompts the user to use the help command.
func main() {
	args := os.Args[1:]
	userRepo := db.NewJSONUserRepo("data/users.json", "data/active_user.json")
	habitRepo := db.NewJSONHabitRepo("data/habits.json")

	if len(args) == 0 {
		fmt.Println("No args! --> Use just <help> entity-arg to see all possible entities!")
		return
	}

	switch args[0] {
	case "user":
		entityArgs := make([]string, len(args[1:]))
		copy(entityArgs, args[1:])
		cli.HandleUser(userRepo, entityArgs)
	case "habit":
		entityArgs := make([]string, len(args[1:]))
		copy(entityArgs, args[1:])
		cli.HandleHabit(habitRepo, userRepo, entityArgs)
	case "help":
		cli.HandleHelp()
	default:
		fmt.Println("Wrong entity! --> Use just <help> entity-arg to see all possible entities!")
	}
}
