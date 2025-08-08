package cli

import (
	"fmt"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
)

// printUsers prints a list of users to the console with their status.
//
//	If the list is empty, it displays a "No users found" message.
//	The title parameter is printed above the list of users.
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

// printHabits prints a list of habits to the console with their status.
//
//	If the list is empty, it displays a "No habits found" message.
//	The title parameter is printed above the list of habits.
func printHabits(habits []domain.Habit, title string) {
	if len(habits) == 0 {
		fmt.Println("No habits found")
		return
	}

	for _, habit := range habits {
		status := "inactive"
		if habit.Status {
			status = "active"
		}

		fmt.Println(title)
		fmt.Printf("Habit ID: %d | User ID: %d | Name: %s "+
			"| Status: %s\n", habit.Id, habit.UserId, habit.Name, status)
	}
}

// HandleUser processes CLI commands related to user management.
//
//	 It delegates execution to the appropriate use case functions based on the provided subcommand in args[0].
//	 Supported commands:
//			active_users, inactive_users, all_users,
//			add,
//			activate, deactivate, update,
//			current, switch,
//			info, all_info.
//
//	 It prints the result or any error messages to the console.
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

// HandleHabit processes CLI commands related to habit management.
//
//	 It delegates execution to the appropriate use case functions based on the provided subcommand in args[0].
//	 Supported commands:
//			add,
//			activate, deact, change,
//			done, undone.
//
//	 It prints the result or any error messages to the console.
func HandleHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) {
	switch args[0] {
	case "add":
		habits, err := usecase.AddHabit(habitRepo, userRepo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printHabits(habits, "== ADD SUCCESS, NEW HABIT ==")
	case "activate":
		habits, err := usecase.ActivateHabit(habitRepo, userRepo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printHabits(habits, "== ACTIVATE SUCCESS, HABIT ==")
	case "deact":
		habits, err := usecase.DeactHabit(habitRepo, userRepo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printHabits(habits, "== DEACTIVATE SUCCESS, HABIT ==")
	case "change":
		habits, err := usecase.ChangeHabit(habitRepo, userRepo, args[1:])
		if err != nil {
			fmt.Println(err)
		}
		printHabits(habits, "== CHANGE SUCCESS, HABIT ==")
	case "done":
		// TODO: usecase.MarkDoneHabit(habitRepo, userRepo, args[1:])
	case "undone":
		// TODO: usecase.MarkUndoneHabit(habitRepo, userRepo, args[1:])
	default:
		fmt.Println("Wrong command! --> Use just <help> entity-arg to see all possible commands!")
	}
}

// HandleHelp prints the CLI help menu to the console.
//
//	It lists available commands for user and habit management,
//	along with a short description of each.
func HandleHelp() {
	fmt.Println("====================================")
	fmt.Println("       Habit Tracker CLI Help       ")
	fmt.Println("====================================")

	fmt.Println("\nUSER COMMANDS:")
	fmt.Printf("  %-35s %s\n", "active_users", "Show all active users")
	fmt.Printf("  %-35s %s\n", "inactive_users", "Show all inactive users")
	fmt.Printf("  %-35s %s\n", "all_users", "Show all users")
	fmt.Printf("  %-35s %s\n", "add <login> <name>", "Add a new user")
	fmt.Printf("  %-35s %s\n", "activate <id>", "Activate a user")
	fmt.Printf("  %-35s %s\n", "deactivate <id>", "Deactivate a user")
	fmt.Printf("  %-35s %s\n", "update <id> <login> <name>", "Update user info")
	fmt.Printf("  %-35s %s\n", "current", "Show current active user")
	fmt.Printf("  %-35s %s\n", "switch <id>", "Switch to another user")
	fmt.Printf("  %-35s %s\n", "info <id>", "Show detailed info about a user (TODO)")
	fmt.Printf("  %-35s %s\n", "all_info", "Show detailed info about all active users (TODO)")

	fmt.Println("\nHABIT COMMANDS:")
	fmt.Printf("  %-35s %s\n", "add <name>", "Add a new habit for the current user")
	fmt.Printf("  %-35s %s\n", "activate <id>", "Activate a habit")
	fmt.Printf("  %-35s %s\n", "deact <id>", "Deactivate a habit")
	fmt.Printf("  %-35s %s\n", "change <id> <name>", "Change habit name")
	fmt.Printf("  %-35s %s\n", "done <id>", "Mark habit as done for today (TODO)")
	fmt.Printf("  %-35s %s\n", "undone <id>", "Mark habit as undone for today (TODO)")

	fmt.Println("\nGENERAL:")
	fmt.Printf("  %-35s %s\n", "help", "Show this help message")
}
