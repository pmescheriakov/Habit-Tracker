package domain

import "time"

type Habit struct {
	Id          int       `json:"Id"`
	UserId      int       `json:"UserId"`
	DateCreated time.Time `json:"DateCreated"`
	Name        string    `json:"Name"`
	Status      bool      `json:"Status"`
}
