package domain

type UserRepository interface {
	GetAll() ([]User, error)
	FindLogName(login, name string) (*User, error)
	FindId(id int) (*User, error)
	Save(user User) error
	Update(user User) error
	Activate(userID int) error
	Deactivate(userID int) error
	SetActive(userID int) error
	GetActive() (*User, error)
}
