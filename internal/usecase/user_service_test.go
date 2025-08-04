package usecase_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
	"github.com/stretchr/testify/require"
)

// setupTestRepo
//
//	Helps to setup temp test repo
func setupTestRepo(t *testing.T) (domain.UserRepository, string, string) {
	tmpDir := t.TempDir()
	usersPath := filepath.Join(tmpDir, "users.json")
	activePath := filepath.Join(tmpDir, "active.json")

	err := os.WriteFile(usersPath, []byte("[]"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(activePath, []byte("[]"), 0644)
	require.NoError(t, err)

	return db.NewJSONUserRepo(usersPath, activePath), usersPath, activePath
}

func TestInactiveUsers(t *testing.T) {

}

func TestActiveUsers(t *testing.T) {

}

func TestAllUsers(t *testing.T) {

}

func TestAddUser(t *testing.T) {

}

func TestActivateUser(t *testing.T) {

}

func TestDeactivateUser(t *testing.T) {

}

func TestUpdateUser(t *testing.T) {

}

func TestCurrentUser(t *testing.T) {

}

func TestSwitchUser(t *testing.T) {

}

func TestInfoUser(t *testing.T) {

}

func TestInfoAllActiveUsers(t *testing.T) {

}
