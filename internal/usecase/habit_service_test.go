package usecase_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
	"github.com/pmescheriakov/Habit-Tracker/internal/usecase"
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

func TestAddHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no args
	habit, err := usecase.AddHabit(habitRepo, userRepo, []string{})
	require.Equal(t, domain.ErrMissingNewHabitFields, err)
	assert.Nil(t, habit)

	// >1 args
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво", "с коллегой"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, habit)

	// no users
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, habit)

	// users exists + active habit exists
	_, err = usecase.AddUser(userRepo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test2", "TestUser2"})
	require.NoError(t, err)

	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.NoError(t, err)
	require.NotNil(t, habit)

	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.Equal(t, domain.ErrHabitExistsActive, err)
	require.Nil(t, habit)

	// users exists + none active habit exists
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})

	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.Equal(t, domain.ErrHabitExistsInactive, err)
	require.Nil(t, habit)

	// users exists + habit doesnt exists
	_, err = usecase.SwitchUser(userRepo, []string{"1"})
	require.NoError(t, err)

	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.NoError(t, err)
	assert.NotNil(t, habit)

	habits, err := habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 2, len(habits))
	assert.Equal(t, "Пить пиво", habits[0].Name)
	assert.Equal(t, 0, habits[1].Id)
	assert.Equal(t, 1, habits[1].UserId)

	// second habit
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Петь караоке"})
	require.NoError(t, err)
	assert.NotNil(t, habit)

	habits, err = habitRepo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 3, len(habits))
	assert.Equal(t, "Пить пиво", habits[0].Name)
	assert.Equal(t, 1, habits[2].Id)
	assert.Equal(t, 1, habits[2].UserId)
}

func TestActivateHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no args
	habit, err := usecase.ActivateHabit(habitRepo, userRepo, []string{})
	require.Equal(t, domain.ErrMissingHabitID, err)
	assert.Nil(t, habit)

	// >1 args
	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"1", "Пить пиво"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, habit)

	// bad 1 arg
	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.Error(t, err)
	assert.Nil(t, habit)

	// no users
	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, habit)

	// no habit exists
	_, err = usecase.AddUser(userRepo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test2", "TestUser2"})
	require.NoError(t, err)

	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habit is active
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.NoError(t, err)

	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrHabitAlreadyActive, err)
	assert.Nil(t, habit)

	// habit is inactive
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.NoError(t, err)

	habit, err = usecase.ActivateHabit(habitRepo, userRepo, []string{"0"})
	require.NoError(t, err)
	assert.NotNil(t, habit)
}

func TestDeactivateHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no args
	habit, err := usecase.DeactivateHabit(habitRepo, userRepo, []string{})
	require.Equal(t, domain.ErrMissingHabitID, err)
	assert.Nil(t, habit)

	// >1 args
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"1", "Пить пиво"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, habit)

	// bad 1 arg
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.Error(t, err)
	assert.Nil(t, habit)

	// no users
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, habit)

	// no habit exists
	_, err = usecase.AddUser(userRepo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test2", "TestUser2"})
	require.NoError(t, err)

	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habit is active
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.NoError(t, err)

	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.NoError(t, err)
	assert.NotNil(t, habit)

	// habit is inactive
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.Equal(t, domain.ErrHabitAlreadyInactive, err)
	assert.Nil(t, habit)
}

func TestUpdateHabit(t *testing.T) {
	habitRepo, userRepo, _, _, _ := setupTestRepos(t)

	// no args
	habit, err := usecase.UpdateHabit(habitRepo, userRepo, []string{})
	require.Equal(t, domain.ErrMissingNewHabitFields, err)
	assert.Nil(t, habit)

	// 1 arg
	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"1"})
	require.Equal(t, domain.ErrMissingNewHabitName, err)
	assert.Nil(t, habit)

	// >2 args
	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"1", "Пить пиво", "с коллегами"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, habit)

	// bad id
	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"Пить пиво", "с коллегами"})
	require.Error(t, err)
	assert.Nil(t, habit)

	// no users
	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"0", "Пить водку"})
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, habit)

	// no habit exists
	_, err = usecase.AddUser(userRepo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	_, err = usecase.AddUser(userRepo, []string{"test2", "TestUser2"})
	require.NoError(t, err)

	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"0", "Пить водку"})
	require.Equal(t, domain.ErrHabitNotFound, err)
	assert.Nil(t, habit)

	// habit is active
	habit, err = usecase.AddHabit(habitRepo, userRepo, []string{"Пить пиво"})
	require.NoError(t, err)

	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"0", "Пить водку"})
	require.NoError(t, err)
	assert.NotNil(t, habit)

	// nothing to change
	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"0", "Пить водку"})
	require.Equal(t, domain.ErrNothingToChange, err)
	assert.Nil(t, habit)

	// habit is inactive
	habit, err = usecase.DeactivateHabit(habitRepo, userRepo, []string{"0"})
	require.NoError(t, err)

	habit, err = usecase.UpdateHabit(habitRepo, userRepo, []string{"0", "Пить водку"})
	require.Equal(t, domain.ErrHabitStatusInactive, err)
	assert.Nil(t, habit)
}

func TestMarkDone(t *testing.T) {
	// TODO: TestMarkDone
}

func TestMarkUndone(t *testing.T) {
	// TODO: TestMarkUndone
}
