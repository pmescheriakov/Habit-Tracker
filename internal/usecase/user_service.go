package usecase

import (
	"errors"
	"strconv"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

func ActiveUsers(userRepo domain.UserRepository) ([]domain.User, error) {
	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	activeUsers := make([]domain.User, 0)
	for _, user := range users {
		if user.Status {
			activeUsers = append(activeUsers, user)
		}
	}

	return activeUsers, err
}

func InactiveUsers(userRepo domain.UserRepository) ([]domain.User, error) {
	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	inactiveUsers := make([]domain.User, 0)
	for _, user := range users {
		if !user.Status {
			inactiveUsers = append(inactiveUsers, user)
		}
	}

	return inactiveUsers, err
}

func AllUsers(userRepo domain.UserRepository) ([]domain.User, error) {
	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, err
}

func AddUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingLoginAndName
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingName
	}
	if len(args) > 2 {
		return nil, domain.ErrTooManyArguments
	}

	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindLogName(args[0], args[1])
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if user != nil && user.Status == true {
		return nil, domain.ErrUserExistsActive
	} else if user != nil {
		return nil, domain.ErrUserExistsInactive
	} else {
		maxId := -1

		for _, u := range users {
			if u.Id > maxId {
				maxId = u.Id
			}
		}

		user = &domain.User{Id: maxId + 1, Login: args[0], Name: args[1], Status: true}

		err = userRepo.Save(*user)
		if err != nil {
			return nil, err
		}

		return append([]domain.User{}, *user), nil
	}
}

func ActivateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingID
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Status == true {
		return append([]domain.User{}, *user), domain.ErrUserAlreadyActive
	}

	err = userRepo.Activate(userId)
	if err != nil {
		return nil, err
	}

	user, err = userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func DeactivateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingID
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Status != true {
		return nil, domain.ErrUserAlreadyInactive
	}

	err = userRepo.Deactivate(userId)
	if err != nil {
		return nil, err
	}

	user, err = userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func UpdateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingNewUserFields
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingNewLoginName
	}
	if len(args) == 2 {
		return nil, domain.ErrMissingNewName
	}
	if len(args) > 3 {
		return nil, domain.ErrTooManyArguments
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Status != true {
		return nil, domain.ErrUserStatusInactive
	}
	if user.Login == args[1] && user.Name == args[2] {
		return nil, domain.ErrNothingToChange
	}

	user.Login = args[1]
	user.Name = args[2]

	err = userRepo.Update(*user)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func CurrentUser(userRepo domain.UserRepository) ([]domain.User, error) {
	user, err := userRepo.GetActive()
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func SwitchUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingActiveUserID
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	userId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	user, err := userRepo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	err = userRepo.SetActive(user.Id)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func InfoUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
	return nil, nil
}

func InfoAllActiveUsers(repo domain.UserRepository, args []string) ([]domain.User, error) {
	return nil, nil
}
