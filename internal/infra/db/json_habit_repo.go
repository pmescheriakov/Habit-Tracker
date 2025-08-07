package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

type JSONHabitRepo struct {
	filePath string
}

func NewJSONHabitRepo(habitsPath string) *JSONHabitRepo {
	_ = os.MkdirAll(filepath.Dir(habitsPath), 0755)

	if _, err := os.Stat(habitsPath); os.IsNotExist(err) {
		_ = os.WriteFile(habitsPath, []byte("[]"), 0644)
	}

	return &JSONHabitRepo{filePath: habitsPath}
}

func (habitsRepo *JSONHabitRepo) writeHabits(habits []domain.Habit) error {
	jsonDb, err := os.OpenFile(habitsRepo.filePath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err = jsonDb.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	data, err := json.MarshalIndent(habits, "", "  ")
	if err != nil {
		return err
	}

	if _, err := jsonDb.Write(data); err != nil {
		return err
	}

	return nil
}

func (habitsRepo *JSONHabitRepo) GetAll() ([]domain.Habit, error) {
	jsonDb, err := os.OpenFile(habitsRepo.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func(jsonDb *os.File) {
		err := jsonDb.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonDb)

	habitsData, err := io.ReadAll(jsonDb)
	if err != nil {
		return nil, err
	}
	if len(habitsData) == 0 {
		return []domain.Habit{}, nil
	}

	var habits []domain.Habit
	err = json.Unmarshal(habitsData, &habits)
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (habitsRepo *JSONHabitRepo) FindName(userId int, name string) (*domain.Habit, error) {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, habit := range habits {
		if habit.Name == name && habit.UserId == userId {
			return &habit, nil
		}
	}

	return nil, domain.ErrHabitNotFound
}

func (habitsRepo *JSONHabitRepo) FindId(userId int, id int) (*domain.Habit, error) {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, habit := range habits {
		if habit.Id == id && habit.UserId == userId {
			return &habit, nil
		}
	}

	return nil, domain.ErrHabitNotFound
}

func (habitsRepo *JSONHabitRepo) Save(habit domain.Habit) error {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return err
	}

	habits = append(habits, habit)

	return habitsRepo.writeHabits(habits)
}

func (habitsRepo *JSONHabitRepo) Update(habit *domain.Habit) error {
	return nil
}

func (habitsRepo *JSONHabitRepo) Activate(id int) error {
	return nil
}

func (habitsRepo *JSONHabitRepo) Deactivate(id int) error {
	return nil
}

func (habitsRepo *JSONHabitRepo) Done(id int) (*domain.Habit, error) {
	return nil, nil
}

func (habitsRepo *JSONHabitRepo) Undone(id int) (*domain.Habit, error) {
	return nil, nil
}
