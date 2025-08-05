package usecase

import (
	"errors"
	"strconv"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

func ActiveUsers(repo domain.UserRepository) ([]domain.User, error) {
	users, err := repo.GetAll()
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

func InactiveUsers(repo domain.UserRepository) ([]domain.User, error) {
	users, err := repo.GetAll()
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

func AllUsers(repo domain.UserRepository) ([]domain.User, error) {
	users, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, err
}

func AddUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingLoginAndName
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingName
	}
	if len(args) > 2 {
		return nil, domain.ErrTooManyArguments
	}

	users, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	user, err := repo.FindLogName(args[0], args[1])
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if user != nil && user.Status == true {
		return nil, domain.ErrUserExistsActive
	} else if user != nil {
		return nil, domain.ErrUserExistsInactive
	} else {
		maxId := -1

		for _, user := range users {
			if user.Id > maxId {
				maxId = user.Id
			}
		}

		user = &domain.User{Id: maxId + 1, Login: args[0], Name: args[1], Status: true}

		err = repo.Save(*user)
		if err != nil {
			return nil, err
		}

		return append([]domain.User{}, *user), nil
	}
}

func ActivateUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
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

	user, err := repo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Status == true {
		return append([]domain.User{}, *user), domain.ErrUserAlreadyActive
	}

	err = repo.Activate(userId)
	if err != nil {
		return nil, err
	}

	user, err = repo.FindId(userId)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func DeactivateUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
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

	user, err := repo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	if user.Status != true {
		return nil, domain.ErrUserAlreadyInactive
	}

	err = repo.Deactivate(userId)
	if err != nil {
		return nil, err
	}

	user, err = repo.FindId(userId)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func UpdateUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
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

	user, err := repo.FindId(userId)
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

	err = repo.Update(*user)
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func CurrentUser(repo domain.UserRepository) ([]domain.User, error) {
	user, err := repo.GetActive()
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

func SwitchUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
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

	user, err := repo.FindId(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	err = repo.SetActive(user.Id)
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
