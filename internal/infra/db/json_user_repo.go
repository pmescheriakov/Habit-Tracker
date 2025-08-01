package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
)

type JSONUserRepo struct {
	filePath string
}

func NewJSONUserRepo(path string) *JSONUserRepo {
	return &JSONUserRepo{filePath: path}
}

func (repo *JSONUserRepo) writeUsers(users []domain.User) error {
	jsonDb, err := os.OpenFile(repo.filePath, os.O_RDWR|os.O_TRUNC, 0666)
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

func (repo *JSONUserRepo) GetAll() ([]domain.User, error) {
	jsonDb, err := os.OpenFile(repo.filePath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
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

func (repo *JSONUserRepo) FindLogName(login, name string) (*domain.User, error) {
	users, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Login == login && user.Name == name {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *JSONUserRepo) FindId(id int) (*domain.User, error) {
	users, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.Id == id {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *JSONUserRepo) Save(user domain.User) error {
	users, err := repo.GetAll()
	if err != nil {
		return err
	}

	users = append(users, user)

	return repo.writeUsers(users)
}

func (repo *JSONUserRepo) Update(user domain.User) error {
	users, err := repo.GetAll()
	if err != nil {
		return err
	}

	users[user.Id] = user

	return repo.writeUsers(users)
}

func (repo *JSONUserRepo) Activate(userID int) error {
	users, err := repo.GetAll()
	if err != nil {
		return err
	}

	users[userID].Status = true

	return repo.writeUsers(users)
}

func (repo *JSONUserRepo) Deactivate(userID int) error {
	users, err := repo.GetAll()
	if err != nil {
		return err
	}

	users[userID].Status = false

	return repo.writeUsers(users)
}

func (repo *JSONUserRepo) SetActive(userID int) error {
	return nil
}

func (repo *JSONUserRepo) GetActive() (*domain.User, error) {
	return nil, nil
}
