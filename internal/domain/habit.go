package domain

import "time"

type Habit struct {
	Id          int
	DateCreated time.Time
	Name        string
	Done        bool
	Active      bool
}
