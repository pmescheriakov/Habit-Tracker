package cli

import (
	"fmt"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
)

func HandleUser(repo domain.UserRepository, args []string) {
	switch args[0] {
	case "users":
		usecase.Users(repo)
	case "add":
		usecase.AddUser(repo, args[1:])
	case "deact":
		usecase.DeactUser(repo, args[1:])
	case "change":
		usecase.ChangeUser(repo, args[1:])
	case "switch":
		usecase.SwitchUser(repo, args[1:])
	case "info":
		usecase.InfoUser(repo, args[1:])
	default:
		fmt.Println("Wrong command! --> Use just <help> entity-arg to see all possible commands!")
	}
}

func HandleHabit(repo domain.UserRepository, args []string) {
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
