package cli

import (
	"fmt"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
)

func HandleUser(args []string) {
	switch args[0] {
	case "users":
		usecase.Users()
	case "add":
		usecase.AddUser(args[1:])
	case "deact":
		usecase.DeactUser(args[1:])
	case "change":
		usecase.ChangeUser(args[1:])
	case "switch":
		usecase.SwitchUser(args[1:])
	case "info":
		usecase.InfoUser(args[1:])
	default:
		fmt.Println("Wrong command! --> Use just <help> entity-arg to see all possible commands!")
	}
}

func HandleHabit(args []string) {
	switch args[0] {
	case "add":
		usecase.AddHabit(args[1:])
	case "stop":
		usecase.StopHabit(args[1:])
	case "change":
		usecase.ChangeHabit(args[1:])
	case "done":
		usecase.MarkDoneHabit(args[1:])
	case "undone":
		usecase.MarkUndoneHabit(args[1:])
	default:
		fmt.Println("Wrong command! --> Use just <help> entity-arg to see all possible commands!")
	}
}

func HandleHelp() {
	// выводить в консоль все типы и команды
}
