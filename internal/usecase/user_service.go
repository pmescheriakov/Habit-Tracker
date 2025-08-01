package usecase

import (
	"fmt"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

func Users(repo domain.UserRepository) {

}

func AddUser(repo domain.UserRepository, args []string) {
	// validate args
	if len(args) == 0 {
		fmt.Println("No <login> <name> provided")
		return
	}
	if len(args) == 1 {
		fmt.Println("No <name> provided")
		return
	}

	users, err := repo.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := repo.Find(args[0], args[1])
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

func DeactUser(repo domain.UserRepository, args []string) {

}

func ChangeUser(repo domain.UserRepository, args []string) {

}

func SwitchUser(repo domain.UserRepository, args []string) {

}

func InfoUser(repo domain.UserRepository, args []string) {

}
