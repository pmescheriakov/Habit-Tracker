package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
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

//func TestGetAllHabits(t *testing.T) {
//	habitRepo, userRepo, _, _, _ := setupTestRepos(t)
//
//	// no habits
//	habits, err := habitRepo.GetAll()
//	require.NoError(t, err)
//	assert.Equal(t, 0, len(habits))
//
//	// some habits
//	err = userRepo.Save()
//
//	err = habitRepo.Save(domain.Habit{Id: 0, UserId: 0, DateCreated: time.Now(), Name: "Пить пиво", Status: true})
//	require.NoError(t, err)
//
//	habits, err = habitRepo.GetAll()
//	require.NoError(t, err)
//	assert.Equal(t, 2, len(habits))
//}
