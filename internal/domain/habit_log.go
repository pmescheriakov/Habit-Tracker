package domain

import "time"

type HabitLog struct {
	Id          int
	DateCreated time.Time
	UserId      int
	HabitId     int
}
