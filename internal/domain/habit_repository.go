package domain

type HabitRepository interface {
	GetAll() ([]Habit, error)
	FindName(userId int, name string) (*Habit, error)
	FindId(userId int, id int) (*Habit, error)
	Save(habit Habit) error
	Update(habit *Habit) error
	Activate(id int) error
	Deactivate(id int) error
	Done(id int) (*Habit, error)
	Undone(id int) (*Habit, error)
}
