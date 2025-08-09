package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

// JSONUserRepo is a JSON-based implementation of the UserRepository interface.
//
//	It stores user data and the active user ID in JSON files at the specified paths.
//	Provides methods for CRUD operations, user activation/deactivation, and managing the active user.
type JSONUserRepo struct {
	filePath   string
	activePath string
}

// NewJSONUserRepo creates a new JSONUserRepo instance.
//
//	It ensures the directories for data files exist and initializes JSON files if they are missing.
//	Returns a pointer to the newly created JSONUserRepo.
func NewJSONUserRepo(usersPath, activeUserPath string) *JSONUserRepo {
	_ = os.MkdirAll(filepath.Dir(usersPath), 0755)
	_ = os.MkdirAll(filepath.Dir(activeUserPath), 0755)

	if _, err := os.Stat(usersPath); os.IsNotExist(err) {
		_ = os.WriteFile(usersPath, []byte("[]"), 0644)
	}
	if _, err := os.Stat(activeUserPath); os.IsNotExist(err) {
		_ = os.WriteFile(activeUserPath, []byte("[]"), 0644)
	}

	return &JSONUserRepo{filePath: usersPath, activePath: activeUserPath}
}

// writeUsers writes the provided list of users to the JSON file.
//
//	It overwrites the entire file with the given slice of users using indentation for readability.
//	Returns an error if writing fails.
func (usersRepo *JSONUserRepo) writeUsers(users []domain.User) error {
	jsonDb, err := os.OpenFile(usersRepo.filePath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err = jsonDb.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	if _, err := jsonDb.Write(data); err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all users from the JSON file.
//
//	Returns a slice of users or an error if reading or unmarshalling fails.
//	If the file is empty, returns an empty slice.
func (usersRepo *JSONUserRepo) GetAll() ([]domain.User, error) {
	jsonDb, err := os.OpenFile(usersRepo.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer func(jsonDb *os.File) {
		err := jsonDb.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonDb)

	usersData, err := io.ReadAll(jsonDb)
	if err != nil {
		return nil, err
	}
	if len(usersData) == 0 {
		return []domain.User{}, nil
	}

	var users []domain.User
	err = json.Unmarshal(usersData, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindLogName searches for a user by login and name.
//
//	Returns a pointer to the user if found, or ErrUserNotFound if not found.
func (usersRepo *JSONUserRepo) FindLogName(login, name string) (*domain.User, error) {
	users, err := usersRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Login == login && user.Name == name {
			return &user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// FindId searches for a user by their ID.
//
//	Returns a pointer to the user if found, or ErrUserNotFound if not found.
func (usersRepo *JSONUserRepo) FindId(id int) (*domain.User, error) {
	users, err := usersRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Id == id {
			return &user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// Save appends a new user to the JSON file.
//
//	Reads all users, appends the new one, and writes back the updated slice.
//	Returns an error if reading or writing fails.
func (usersRepo *JSONUserRepo) Save(user domain.User) error {
	users, err := usersRepo.GetAll()
	if err != nil {
		return err
	}

	users = append(users, user)

	return usersRepo.writeUsers(users)
}

// Update replaces an existing user in the JSON file.
//
//	Searches for the user by ID and replaces it with the provided user data.
//	Returns an error if the user is not found or if reading/writing fails.
func (usersRepo *JSONUserRepo) Update(user domain.User) error {
	u, err := usersRepo.FindId(user.Id)
	if err != nil {
		return err
	}
	if u == nil {
		return domain.ErrUserNotFound
	}

	users, err := usersRepo.GetAll()
	if err != nil {
		return err
	}

	users[user.Id] = user

	return usersRepo.writeUsers(users)
}

// Activate sets the status of a user to active (true) by their ID.
//
//	Returns ErrUserNotFound if the user does not exist.
//	Returns an error if reading or writing fails.
func (usersRepo *JSONUserRepo) Activate(userID int) error {
	users, err := usersRepo.GetAll()
	if err != nil {
		return err
	}

	checkUser, err := usersRepo.FindId(userID)
	if err != nil {
		return err
	}
	if checkUser == nil {
		return domain.ErrUserNotFound
	}

	users[userID].Status = true

	return usersRepo.writeUsers(users)
}

// Deactivate sets the status of a user to inactive (false) by their ID.
//
//	Returns ErrUserNotFound if the user does not exist.
//	Returns an error if reading or writing fails.
func (usersRepo *JSONUserRepo) Deactivate(userID int) error {
	users, err := usersRepo.GetAll()
	if err != nil {
		return err
	}

	checkUser, err := usersRepo.FindId(userID)
	if err != nil {
		return err
	}
	if checkUser == nil {
		return domain.ErrUserNotFound
	}

	users[userID].Status = false

	return usersRepo.writeUsers(users)
}

// SetNewCurrent sets the active user ID in the active user JSON file.
//
//	Validates that the user exists before writing the active user ID.
//	Returns an error if the user does not exist or if writing fails.
func (usersRepo *JSONUserRepo) SetNewCurrent(userID int) error {
	jsonDb, err := os.OpenFile(usersRepo.activePath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(jsonDb *os.File) {
		err := jsonDb.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonDb)

	checkUser, err := usersRepo.FindId(userID)
	if err != nil {
		return err
	}
	if checkUser == nil {
		return domain.ErrUserNotFound
	}

	userActive := struct {
		ActiveUserId int `json:"activeUserId"`
	}{
		ActiveUserId: userID,
	}

	data, err := json.MarshalIndent(userActive, "", "  ")
	if err != nil {
		return err
	}

	if _, err := jsonDb.Write(data); err != nil {
		return err
	}

	return nil
}

// GetCurrent retrieves the currently active user from the JSON file.
//
//	If the active user ID is not set, initializes it to 0.
//	Returns a pointer to the active user or an error if retrieval fails.
func (usersRepo *JSONUserRepo) GetCurrent() (*domain.User, error) {
	jsonDb, err := os.OpenFile(usersRepo.activePath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	defer func(jsonDb *os.File) {
		err := jsonDb.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(jsonDb)

	bytes, err := io.ReadAll(jsonDb)
	if err != nil {
		return nil, err
	}

	users, err := usersRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, domain.ErrNoUsers
	}

	var userActive struct {
		ActiveUserId int `json:"activeUserId"`
	}

	err = json.Unmarshal(bytes, &userActive)
	if err != nil {
		// if no active user -> set 0
		userActive.ActiveUserId = 0

		data, err := json.MarshalIndent(userActive, "", "  ")
		if err != nil {
			return nil, err
		}

		if _, err := jsonDb.Seek(0, 0); err != nil {
			return nil, err
		}
		if err := jsonDb.Truncate(0); err != nil {
			return nil, err
		}
		if _, err := jsonDb.Write(data); err != nil {
			return nil, err
		}
	}

	return &users[userActive.ActiveUserId], nil
}
