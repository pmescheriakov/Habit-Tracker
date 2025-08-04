package cli

import (
	"fmt"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
)

func printUsers(users []domain.User, title string) {
	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	for _, user := range users {
		status := "inactive"
		if user.Status {
			status = "active"
		}

		fmt.Println(title)
		fmt.Printf("ID: %d | Login: %s | Name: %s | Status: %s\n", user.Id, user.Login, user.Name, status)
	}
}

func HandleUser(repo domain.UserRepository, args []string) {
	switch args[0] {
	case "active_users":
		users, err := usecase.ActiveUsers(repo)
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== ALL ACTIVE USERS ==")
	case "inactive_users":
		users, err := usecase.InactiveUsers(repo)
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== ALL INACTIVE USERS ==")
	case "all_users":
		users, err := usecase.AllUsers(repo)
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== ALL USERS ==")
	case "add":
		users, err := usecase.AddUser(repo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== ADD SUCCESS, NEW USER ==")
	case "activate":
		users, err := usecase.ActivateUser(repo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== ACTIVATE SUCCESS, ACTIVATED USER ==")
	case "deactivate":
		users, err := usecase.DeactivateUser(repo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== DEACTIVATE SUCCESS, DEACTIVATED USER ==")
	case "update":
		users, err := usecase.UpdateUser(repo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== UPDATE SUCCESS, UPDATED USER ==")
	case "current":
		users, err := usecase.CurrentUser(repo)
		if err != nil {
			fmt.Println(err)
		}
		printUsers(users, "== CURRENT USER ==")
	case "switch":
		users, err := usecase.SwitchUser(repo, args[1:])
		if err != nil {
			fmt.Println(err)
		}

		users, err = usecase.CurrentUser(repo)
		if err != nil {
			fmt.Println(err)
		}

		printUsers(users, "== CURRENT USER ==")
	case "info":
		// TODO: usecase.InfoUser(repo, args[1:])
	case "all_info":
		// TODO: usecase.InfoAllActiveUsers(repo, args[1:])
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
