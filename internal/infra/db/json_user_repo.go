package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

type JSONUserRepo struct {
	filePath   string
	activePath string
}

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

func (usersRepo *JSONUserRepo) Save(user domain.User) error {
	users, err := usersRepo.GetAll()
	if err != nil {
		return err
	}

	users = append(users, user)

	return usersRepo.writeUsers(users)
}

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

func (usersRepo *JSONUserRepo) SetActive(userID int) error {
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

func (usersRepo *JSONUserRepo) GetActive() (*domain.User, error) {
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

		if _, err := jsonDb.Write(data); err != nil {
			return nil, err
		}
	}

	return &users[userActive.ActiveUserId], nil
}
