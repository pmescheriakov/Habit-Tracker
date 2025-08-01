package domain

type UserRepository interface {
	GetAll() ([]User, error)
	Find(login, name string) (*User, error)
	Save(user User) error
	Update(user User) error
	Activate(userID int) error
	Deactivate(userID int) error
	SetActive(userID int) error
	GetActive() (*User, error)
}
