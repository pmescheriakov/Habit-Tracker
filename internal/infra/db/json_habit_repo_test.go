package db_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepos
//
//	Helps to setup temp test repos
func setupTestRepos(t *testing.T) (domain.HabitRepository, domain.UserRepository, string, string, string) {
	tmpDir := t.TempDir()
	usersPath := filepath.Join(tmpDir, "users.json")
	activePath := filepath.Join(tmpDir, "active.json")
	habitsPath := filepath.Join(tmpDir, "habits.json")

	err := os.WriteFile(usersPath, []byte("[]"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(activePath, []byte("[]"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(habitsPath, []byte("[]"), 0644)
	require.NoError(t, err)

	return db.NewJSONHabitRepo(habitsPath), db.NewJSONUserRepo(usersPath, activePath), habitsPath, usersPath, activePath
}

func TestGetAllHabits(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no habits
	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 0, len(habits))

	// some habits and user exists
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = userRepo.Save(user0)
	require.NoError(t, err)

	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	habit1 := domain.Habit{Id: 1, UserId: 0, DateCreated: time.Now(), Name: "Таблетки", Status: true}
	err = habitRepo.Save(habit1)
	require.NoError(t, err)

	habits, err = habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 2, len(habits))

	names := []string{habits[0].Name, habits[1].Name}

	expectedNames := map[string]bool{
		"Пить пиво": true,
		"Таблетки":  true,
	}

	for _, name := range names {
		assert.True(t, expectedNames[name])
	}
}

func TestFindName(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no user
	habit, err := habitRepo.FindName(1, "Drink vodka")
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// no habits on user
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = userRepo.Save(user0)
	require.NoError(t, err)

	habit, err = habitRepo.FindName(0, "Drink vodka")
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habits + not exists
	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	habit1 := domain.Habit{Id: 1, UserId: 0, DateCreated: time.Now(), Name: "Таблетки", Status: true}
	err = habitRepo.Save(habit1)
	require.NoError(t, err)

	habit, err = habitRepo.FindName(0, "Drink vodka")
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habits + exists
	habit, err = habitRepo.FindName(0, "Пить пиво")
	require.NoError(t, err)
	assert.NotNil(t, habit)
}

func TestFindId(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no user
	habit, err := habitRepo.FindId(1, 2)
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// no habits on user
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = userRepo.Save(user0)
	require.NoError(t, err)

	habit, err = habitRepo.FindId(0, 2)
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habits + not exists
	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	habit1 := domain.Habit{Id: 1, UserId: 0, DateCreated: time.Now(), Name: "Таблетки", Status: true}
	err = habitRepo.Save(habit1)
	require.NoError(t, err)

	habit, err = habitRepo.FindId(0, 2)
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habits + exists
	habit, err = habitRepo.FindId(0, 0)
	require.NoError(t, err)
	assert.NotNil(t, habit)
}

func TestSaveHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err := userRepo.Save(user0)
	require.NoError(t, err)

	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(habits))
	assert.Equal(t, habit0.Name, habits[0].Name)
	assert.Equal(t,
		habit0.DateCreated.Truncate(time.Second),
		habits[0].DateCreated.Truncate(time.Second),
	)
}

func TestUpdateHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err := userRepo.Save(user0)
	require.NoError(t, err)

	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	habit0New := habit0
	habit0New.Name = "Drink vodka"

	err = habitRepo.Update(habit0New, 0)
	require.NoError(t, err)

	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(habits))
	assert.Equal(t, habit0New.Name, habits[0].Name)
	assert.Equal(t,
		habit0.DateCreated.Truncate(time.Second),
		habits[0].DateCreated.Truncate(time.Second),
	)
}

func TestActivate(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err := userRepo.Save(user0)
	require.NoError(t, err)

	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: false}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	err = habitRepo.Activate(0, 0)
	require.NoError(t, err)

	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, true, habits[0].Status)
}

func TestDeactivate(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err := userRepo.Save(user0)
	require.NoError(t, err)

	habit0 := domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true}
	err = habitRepo.Save(habit0)
	require.NoError(t, err)

	err = habitRepo.Deactivate(0, 0)
	require.NoError(t, err)

	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, false, habits[0].Status)
}
