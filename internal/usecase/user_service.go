package usecase

import (
	"fmt"
	"strconv"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

func Users(repo domain.UserRepository) {

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
	if len(args) > 1 {
		fmt.Println("Too many arguments")
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
		fmt.Println("User not exists")

		user = &domain.User{Id: len(users), Login: args[0], Name: args[1], Status: true}

		err = repo.Save(*user)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("User added! Id: %v, Login: %v, Name: %v\n", user.Id, user.Login, user.Name)
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

	fmt.Printf("User activated! Id: %v, Login: %v, Name: %v\n", user.Id, user.Login, user.Name)
}

func DeactUser(repo domain.UserRepository, args []string) {
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

	fmt.Printf("User deactivated! Id: %v, Login: %v, Name: %v\n", user.Id, user.Login, user.Name)
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

	fmt.Printf("User updated! Id: %v, Login: %v, Name: %v\n", user.Id, user.Login, user.Name)
}

func SwitchUser(repo domain.UserRepository, args []string) {

}

func InfoUser(repo domain.UserRepository, args []string) {

}
