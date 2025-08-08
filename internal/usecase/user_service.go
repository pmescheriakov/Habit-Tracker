package usecase

import (
	"errors"
	"strconv"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

// ActiveUsers retrieves all active users from the repository.
//
//	It filters the complete user list, returning only those with Status set to true.
//	Returns the list of active users or an error if retrieval fails.
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

// InactiveUsers retrieves all inactive users from the repository.
//
//	It filters the complete user list, returning only those with Status set to false.
//	Returns the list of inactive users or an error if retrieval fails.
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

// AllUsers retrieves all users from the repository.
//
//	Returns the list of all users or an error if retrieval fails.
func AllUsers(userRepo domain.UserRepository) ([]domain.User, error) {
	users, err := userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return users, err
}

// AddUser creates a new active user with the provided login and name.
//
//	It validates the input arguments, ensures no active user with the same login and name exists,
//	and assigns a unique ID to the new user.
//	Returns the created user or an error if the operation fails.
func AddUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingLoginAndName
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingUserName
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

// ActivateUser sets a user's status to active based on the provided user ID.
//
//	It validates the ID, ensures the user exists and is currently inactive,
//	then updates the status to active.
//	Returns the updated user or an error if the operation fails.
func ActivateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingUserID
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

// DeactivateUser sets a user's status to inactive based on the provided user ID.
//
//	It validates the ID, ensures the user exists and is currently active,
//	then updates the status to inactive.
//	Returns the updated user or an error if the operation fails.
func DeactivateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingUserID
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

// UpdateUser changes the login and name of an active user.
//
//	It validates the provided arguments, ensures the user exists, is active,
//	and that at least one of the fields is different before updating.
//	Returns the updated user or an error if the operation fails.
func UpdateUser(userRepo domain.UserRepository, args []string) ([]domain.User, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingNewUserFields
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingNewUserLoginName
	}
	if len(args) == 2 {
		return nil, domain.ErrMissingNewUserName
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

// CurrentUser retrieves the currently active user.
//
//	Returns a slice containing the active user or an error if retrieval fails.
func CurrentUser(userRepo domain.UserRepository) ([]domain.User, error) {
	user, err := userRepo.GetActive()
	if err != nil {
		return nil, err
	}

	return append([]domain.User{}, *user), nil
}

// SwitchUser changes the active user to the one specified by ID.
//
//	It validates the ID, ensures the user exists,
//	and sets the specified user as the current active user.
//	Returns the newly active user or an error if the operation fails.
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

// InfoUser retrieves detailed information about a specific user.
//
//	Currently a placeholder for future implementation.
func InfoUser(repo domain.UserRepository, args []string) ([]domain.User, error) {
	return nil, nil
}

// InfoAllActiveUsers retrieves detailed information about all active users.
//
//	Currently a placeholder for future implementation.
func InfoAllActiveUsers(repo domain.UserRepository, args []string) ([]domain.User, error) {
	return nil, nil
}
