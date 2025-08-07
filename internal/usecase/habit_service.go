package usecase

import (
	"errors"
	"time"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

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

		for _, h := range habits {
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

func ActivateHabit(repo domain.HabitRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}

func DeactHabit(repo domain.HabitRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}

func ChangeHabit(repo domain.HabitRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}

func MarkDoneHabit(repo domain.HabitRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}

func MarkUndoneHabit(repo domain.HabitRepository, args []string) ([]domain.Habit, error) {
	return nil, nil
}
