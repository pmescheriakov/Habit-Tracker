package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

// JSONHabitRepo is a JSON-based implementation of the HabitRepository interface.
//
//	It stores habit data in a JSON file at the specified filePath.
//	Provides methods for CRUD operations and habit status changes.
type JSONHabitRepo struct {
	filePath string
}

// NewJSONHabitRepo creates a new JSONHabitRepo instance.
//
//	It ensures the directory exists and initializes the JSON file if it does not exist.
//	Returns a pointer to the newly created JSONHabitRepo.
func NewJSONHabitRepo(habitsPath string) *JSONHabitRepo {
	_ = os.MkdirAll(filepath.Dir(habitsPath), 0755)

	if _, err := os.Stat(habitsPath); os.IsNotExist(err) {
		_ = os.WriteFile(habitsPath, []byte("[]"), 0644)
	}

	return &JSONHabitRepo{filePath: habitsPath}
}

// writeHabits writes the provided list of habits to the JSON file.
//
//	It overwrites the entire file with the given habits slice using indentation for readability.
//	Returns an error if writing fails.
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

// GetAll retrieves all habits from the JSON file.
//
//	Returns a slice of habits or an error if reading or unmarshalling fails.
//	If the file is empty, returns an empty slice.
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

// FindName searches for a habit by user ID and habit name.
//
//	Returns a pointer to the habit if found, or ErrHabitNotFound if not found.
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

// FindId searches for a habit by user ID and habit ID.
//
//	Returns a pointer to the habit if found, or ErrHabitNotFound if not found.
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

// Save appends a new habit to the JSON file.
//
//	Reads all habits, appends the new habit, and writes back the updated slice.
//	Returns an error if reading or writing fails.
func (habitsRepo *JSONHabitRepo) Save(habit domain.Habit) error {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return err
	}

	habits = append(habits, habit)

	return habitsRepo.writeHabits(habits)
}

// Update replaces an existing habit for the given user ID.
//
//	Searches for the habit by ID and user ID, replaces it with the provided habit,
//	and writes the updated slice back to the file.
//	Returns an error if reading or writing fails.
func (habitsRepo *JSONHabitRepo) Update(habit domain.Habit, userId int) error {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return err
	}

	for i := range habits {
		if habits[i].UserId == userId && habits[i].Id == habit.Id {
			habits[i] = habit
		}
	}

	return habitsRepo.writeHabits(habits)
}

// Activate sets the status of a habit to active (true) for a given user ID and habit ID.
//
//	Returns ErrHabitNotFound if the habit does not exist.
//	Returns an error if reading or writing fails.
func (habitsRepo *JSONHabitRepo) Activate(id int, userId int) error {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return err
	}

	checkHabit, err := habitsRepo.FindId(userId, id)
	if err != nil {
		return err
	}
	if checkHabit == nil {
		return domain.ErrHabitNotFound
	}

	for i := range habits {
		if habits[i].UserId == userId && habits[i].Id == id {
			habits[i].Status = true
		}
	}

	return habitsRepo.writeHabits(habits)
}

// Deactivate sets the status of a habit to inactive (false) for a given user ID and habit ID.
//
//	Returns ErrHabitNotFound if the habit does not exist.
//	Returns an error if reading or writing fails.
func (habitsRepo *JSONHabitRepo) Deactivate(id int, userId int) error {
	habits, err := habitsRepo.GetAll()
	if err != nil {
		return err
	}

	checkHabit, err := habitsRepo.FindId(userId, id)
	if err != nil {
		return err
	}
	if checkHabit == nil {
		return domain.ErrHabitNotFound
	}

	for i := range habits {
		if habits[i].UserId == userId && habits[i].Id == id {
			habits[i].Status = false
		}
	}

	return habitsRepo.writeHabits(habits)
}

// Done marks a habit as done for the current day.
//
//	Currently a placeholder for future implementation.
//	Should update the habit log to reflect today's completion.
func (habitsRepo *JSONHabitRepo) Done(id int) (*domain.Habit, error) {
	return nil, nil
}

// Undone marks a habit as undone for the current day.
//
//	Currently a placeholder for future implementation.
//	Should update the habit log to remove today's completion.
func (habitsRepo *JSONHabitRepo) Undone(id int) (*domain.Habit, error) {
	return nil, nil
}
