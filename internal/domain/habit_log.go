package domain

import "time"

type HabitLog struct {
	Hab        Habit
	Done       bool
	DateLogged time.Time
}
