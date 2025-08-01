package cli

import (
	"fmt"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
)

func HandleUser(repo domain.UserRepository, args []string) {
	switch args[0] {
	case "active_users":
		usecase.ActiveUsers(repo)
	case "inactive_users":
		usecase.InactiveUsers(repo)
	case "all_users":
		usecase.AllUsers(repo)
	case "add":
		usecase.AddUser(repo, args[1:])
	case "activate":
		usecase.ActivateUser(repo, args[1:])
	case "deactivate":
		usecase.DeactivateUser(repo, args[1:])
	case "update":
		usecase.UpdateUser(repo, args[1:])
	case "current":
		usecase.CurrentUser(repo)
	case "switch":
		usecase.SwitchUser(repo, args[1:])
	case "info":
		usecase.InfoUser(repo, args[1:])
	case "all_info":
		usecase.InfoAllActiveUsers(repo, args[1:])
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
