package usecase

import (
	"errors"
	"strconv"
	"time"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

// AddHabit adds a new habit for the current active user.
//
//	It validates the provided arguments, checks if a habit with the same name exists,
//	and creates a new active habit if no conflicts are found.
//	Returns the newly created habit or an error if the operation fails.
func AddHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingNewHabitFields
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	habits, err := habitRepo.GetAll()
	if err != nil {
		return nil, err
	}

	currentUser, err := CurrentUser(userRepo)
	if err != nil {
		return nil, err
	}

	habit, err := habitRepo.FindName(currentUser[0].Id, args[0])
	if err != nil && !errors.Is(err, domain.ErrHabitNotFound) {
		return nil, err
	}

	if habit != nil && habit.Status == true {
		return nil, domain.ErrHabitExistsActive
	} else if habit != nil {
		return nil, domain.ErrHabitExistsInactive
	} else {
		maxId := -1

		userHabits := make([]domain.Habit, 0)

		for _, habit := range habits {
			if habit.UserId == currentUser[0].Id && habit.Id > maxId {
				userHabits = append(userHabits, habit)
			}
		}

		for _, h := range userHabits {
			if h.Id > maxId {
				maxId = h.Id
			}
		}

		habit = &domain.Habit{Id: maxId + 1, UserId: currentUser[0].Id, DateCreated: time.Now(), Name: args[0], Status: true}

		err = habitRepo.Save(*habit)
		if err != nil {
			return nil, err
		}

		return append([]domain.Habit{}, *habit), nil
	}
}

// ActivateHabit activates an inactive habit for the current active user.
//
//	It validates the provided habit ID, ensures the habit exists and is inactive,
//	and updates its status to active.
//	Returns the updated habit or an error if the operation fails.
func ActivateHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingHabitID
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	habitId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	currentUser, err := CurrentUser(userRepo)
	if err != nil {
		return nil, err
	}

	habit, err := habitRepo.FindId(currentUser[0].Id, habitId)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, domain.ErrHabitNotFound
	}
	if habit.Status == true {
		return nil, domain.ErrHabitAlreadyActive
	}

	err = habitRepo.Activate(habitId, currentUser[0].Id)
	if err != nil {
		return nil, err
	}

	habit, err = habitRepo.FindId(currentUser[0].Id, habitId)
	if err != nil {
		return nil, err
	}

	return append([]domain.Habit{}, *habit), nil
}

// DeactivateHabit deactivates an active habit for the current active user.
//
//	It validates the provided habit ID, ensures the habit exists and is active,
//	and updates its status to inactive.
//	Returns the updated habit or an error if the operation fails.
func DeactivateHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingHabitID
	}
	if len(args) > 1 {
		return nil, domain.ErrTooManyArguments
	}

	habitId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	currentUser, err := CurrentUser(userRepo)
	if err != nil {
		return nil, err
	}

	habit, err := habitRepo.FindId(currentUser[0].Id, habitId)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, domain.ErrHabitNotFound
	}
	if habit.Status != true {
		return nil, domain.ErrHabitAlreadyInactive
	}

	err = habitRepo.Deactivate(habitId, currentUser[0].Id)
	if err != nil {
		return nil, err
	}

	habit, err = habitRepo.FindId(currentUser[0].Id, habitId)
	if err != nil {
		return nil, err
	}

	return append([]domain.Habit{}, *habit), nil
}

// UpdateHabit updates the name of an active habit for the current active user.
//
//	It validates the provided habit ID and new name, ensures the habit exists,
//	is active, and the name is different, then updates it.
//	Returns the updated habit or an error if the operation fails.
func UpdateHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	if len(args) == 0 {
		return nil, domain.ErrMissingNewHabitFields
	}
	if len(args) == 1 {
		return nil, domain.ErrMissingNewHabitName
	}
	if len(args) > 2 {
		return nil, domain.ErrTooManyArguments
	}

	habitId, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	currentUser, err := CurrentUser(userRepo)
	if err != nil {
		return nil, err
	}

	habit, err := habitRepo.FindId(currentUser[0].Id, habitId)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, domain.ErrHabitNotFound
	}
	if habit.Status != true {
		return nil, domain.ErrHabitStatusInactive
	}
	if habit.Name == args[1] {
		return nil, domain.ErrNothingToChange
	}

	habit.Name = args[1]

	err = habitRepo.Update(*habit, currentUser[0].Id)
	if err != nil {
		return nil, err
	}

	return append([]domain.Habit{}, *habit), nil
}

// MarkDoneHabit marks a habit as done for the current day.
//
//	Currently a placeholder for implementation.
//	Should update the habit log to reflect today's completion status.
func MarkDoneHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}

// MarkUndoneHabit marks a habit as undone for the current day.
//
//	Currently a placeholder for implementation.
//	Should update the habit log to remove today's completion status.
func MarkUndoneHabit(habitRepo domain.HabitRepository, userRepo domain.UserRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}
