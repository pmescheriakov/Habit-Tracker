package usecase

import (
	"fmt"
	"strconv"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

func printUsers(users []domain.User) {
	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}

	for _, user := range users {
		status := "inactive"
		if user.Status {
			status = "active"
		}

		fmt.Printf("ID: %d | Login: %s | Name: %s | Status: %s\n", user.Id, user.Login, user.Name, status)
	}
}

func InactiveUsers(repo domain.UserRepository) {
	users, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	activeUsers := make([]domain.User, 0)
	for _, user := range users {
		if !user.Status {
			activeUsers = append(activeUsers, user)
		}
	}

	printUsers(activeUsers)
}

func ActiveUsers(repo domain.UserRepository) {
	users, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	activeUsers := make([]domain.User, 0)
	for _, user := range users {
		if user.Status {
			activeUsers = append(activeUsers, user)
		}
	}

	printUsers(activeUsers)
}

func AllUsers(repo domain.UserRepository) {
	users, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	printUsers(users)
}

func AddUser(repo domain.UserRepository, args []string) {
	if len(args) == 0 {
		fmt.Println("No <login> <name> provided")
		return
	}
	if len(args) == 1 {
		fmt.Println("No <name> provided")
		return
	}
	if len(args) > 2 {
		fmt.Println("Too many arguments")
		return
	}

	users, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.FindLogName(args[0], args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	if user != nil && user.Status == true {
		fmt.Println("User already exists and active")
	} else if user != nil {
		fmt.Println("User already exists and inactive")
	} else {
		user = &domain.User{Id: len(users), Login: args[0], Name: args[1], Status: true}

		err = repo.Save(*user)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("User added!")
		printUsers(append([]domain.User{}, *user))
	}
}

func ActivateUser(repo domain.UserRepository, args []string) {
	if len(args) == 0 {
		fmt.Println("No user <id> provided")
		return
	}
	if len(args) > 1 {
		fmt.Println("Too many arguments")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.FindId(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}
	if user.Status == true {
		fmt.Println("User already active")
		return
	}

	err = repo.Activate(userId)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User activated!")
	printUsers(append([]domain.User{}, *user))
}

func DeactivateUser(repo domain.UserRepository, args []string) {
	if len(args) == 0 {
		fmt.Println("No user <id> provided")
		return
	}
	if len(args) > 1 {
		fmt.Println("Too many arguments")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.FindId(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}
	if user.Status != true {
		fmt.Println("User already inactive")
		return
	}

	err = repo.Deactivate(userId)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User deactivated!")
	printUsers(append([]domain.User{}, *user))
}

func UpdateUser(repo domain.UserRepository, args []string) {
	if len(args) == 0 {
		fmt.Println("No user <id>, new <login>, new <name> provided")
		return
	}
	if len(args) == 1 {
		fmt.Println("No user new <login>, new <name> provided")
	}
	if len(args) == 2 {
		fmt.Println("No user new <name> provided")
	}
	if len(args) > 3 {
		fmt.Println("Too many arguments")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.FindId(userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}
	if user.Status != true {
		fmt.Println("User inactive. Activate user before update!")
		return
	}
	if user.Name == args[1] && user.Login == args[2] {
		fmt.Println("Nothing to change")
		return
	}

	user.Name = args[1]
	user.Login = args[2]

	err = repo.Update(*user)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User updated!")
	printUsers(append([]domain.User{}, *user))
}

func CurrentUser(repo domain.UserRepository) {
	user, err := repo.GetActive()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Active User Info")
	printUsers(append([]domain.User{}, *user))
}

func SwitchUser(repo domain.UserRepository, args []string) {
	if len(args) == 0 {
		fmt.Println("No new active user <id> provided")
		return
	}
	if len(args) > 1 {
		fmt.Println("Too many arguments")
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.FindId(userId)
	if err != nil {
		fmt.Println(err)
	}
	if user == nil {
		fmt.Println("User not found")
		return
	}

	err = repo.SetActive(user.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User updated!")
	printUsers(append([]domain.User{}, *user))

	CurrentUser(repo)
}

func InfoUser(repo domain.UserRepository, args []string) {

}

func InfoAllActiveUsers(repo domain.UserRepository, args []string) {

}
