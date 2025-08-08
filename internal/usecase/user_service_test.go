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

func TestActiveUsers(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no users
	users, err := usecase.ActiveUsers(repo)
	require.NoError(t, err)
	assert.Equal(t, []domain.User{}, users)

	// users
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})
	require.NoError(t, err)

	_, err = usecase.DeactivateUser(repo, []string{"1"})

	users, err = usecase.ActiveUsers(repo)
	require.NoError(t, err)
	require.Equal(t, 2, len(users))
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}, users[0])
	assert.Equal(t, domain.User{Id: 2, Login: "test2", Name: "TestUser2", Status: true}, users[1])
}

func TestInactiveUsers(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no users
	users, err := usecase.ActiveUsers(repo)
	require.NoError(t, err)
	assert.Equal(t, []domain.User{}, users)

	// users
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	_, err = usecase.DeactivateUser(repo, []string{"1"})

	users, err = usecase.InactiveUsers(repo)
	require.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: false}, users[0])
}

func TestAllUsers(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no users
	users, err := usecase.ActiveUsers(repo)
	require.NoError(t, err)
	assert.Equal(t, []domain.User{}, users)

	// users
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	_, err = usecase.DeactivateUser(repo, []string{"1"})

	users, err = usecase.AllUsers(repo)
	require.NoError(t, err)
	assert.Equal(t, 3, len(users))
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}, users[0])
	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: false}, users[1])
	assert.Equal(t, domain.User{Id: 2, Login: "test2", Name: "TestUser2", Status: true}, users[2])
}

func TestAddUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no args
	user, err := usecase.AddUser(repo, []string{})
	assert.Equal(t, domain.ErrMissingLoginAndName, err)
	assert.Nil(t, user)

	// one arg
	user, err = usecase.AddUser(repo, []string{"skebob"})
	assert.Equal(t, domain.ErrMissingUserName, err)
	assert.Nil(t, user)

	// many args
	user, err = usecase.AddUser(repo, []string{"test0", "TestUser0", "skebob"})
	assert.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, user)

	// 2 args
	user, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	require.NoError(t, err)
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}, user[0])

	// try to add same active user
	user, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	require.Equal(t, domain.ErrUserExistsActive, err)
	assert.Nil(t, user)

	// try to add same deact user
	pUser, err := repo.FindLogName("test0", "TestUser0")
	require.NoError(t, err)

	err = repo.Deactivate(pUser.Id)
	require.NoError(t, err)

	user, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	require.Equal(t, domain.ErrUserExistsInactive, err)
	assert.Nil(t, user)

	// check users
	user, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	require.NoError(t, err)
	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}, user[0])

	users, err := usecase.AllUsers(repo)
	require.NoError(t, err)
	require.Equal(t, 2, len(users))
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: false}, users[0])
	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}, users[1])
}

func TestActivateUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no args
	user, err := usecase.ActivateUser(repo, []string{})
	require.Equal(t, domain.ErrMissingUserID, err)
	assert.Nil(t, user)

	// many args
	user, err = usecase.ActivateUser(repo, []string{"test0", "TestUser0", "skebob"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, user)

	// one no int arg
	user, err = usecase.ActivateUser(repo, []string{"skebob"})
	require.Error(t, err)
	assert.Nil(t, user)

	// no user
	user, err = usecase.ActivateUser(repo, []string{"1"})
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)

	// already activated
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	user, err = usecase.ActivateUser(repo, []string{"1"})
	require.Equal(t, domain.ErrUserAlreadyActive, err)
	assert.NotNil(t, user)

	// deactivated to activated
	user, err = usecase.DeactivateUser(repo, []string{"1"})
	require.NoError(t, err)

	user, err = usecase.ActivateUser(repo, []string{"1"})
	require.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}, user[0])
}

func TestDeactivateUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no args
	user, err := usecase.DeactivateUser(repo, []string{})
	require.Equal(t, domain.ErrMissingUserID, err)
	assert.Nil(t, user)

	// many args
	user, err = usecase.DeactivateUser(repo, []string{"test0", "TestUser0", "skebob"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, user)

	// one no int arg
	user, err = usecase.DeactivateUser(repo, []string{"skebob"})
	require.Error(t, err)
	assert.Nil(t, user)

	// no user
	user, err = usecase.DeactivateUser(repo, []string{"1"})
	require.EqualError(t, err, "user not found")
	assert.Nil(t, user)

	// already deactivated
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	user, err = usecase.DeactivateUser(repo, []string{"1"})
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: false}, user[0])

	// already deactivated
	user, err = usecase.DeactivateUser(repo, []string{"1"})
	require.Equal(t, domain.ErrUserAlreadyInactive, err)
	assert.Nil(t, user)
}

func TestUpdateUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no args
	user, err := usecase.UpdateUser(repo, []string{})
	require.Equal(t, domain.ErrMissingNewUserFields, err)
	assert.Nil(t, user)

	// one arg
	user, err = usecase.UpdateUser(repo, []string{"1"})
	require.Equal(t, domain.ErrMissingNewUserLoginName, err)
	assert.Nil(t, user)

	// two args
	user, err = usecase.UpdateUser(repo, []string{"0", "test0"})
	require.Equal(t, domain.ErrMissingNewUserName, err)
	assert.Nil(t, user)

	// many args
	user, err = usecase.UpdateUser(repo, []string{"0", "test0", "TestUser0", "skebob"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, user)

	// three, no int arg
	user, err = usecase.UpdateUser(repo, []string{"o", "test0", "TestUser0"})
	require.Error(t, err)
	assert.Nil(t, user)

	// no user
	user, err = usecase.UpdateUser(repo, []string{"4", "test4", "TestUser4"})
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)

	// user inactive
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	user, err = usecase.DeactivateUser(repo, []string{"1"})
	require.NoError(t, err)

	user, err = usecase.UpdateUser(repo, []string{"1", "test11", "TestUser11"})
	require.Equal(t, domain.ErrUserStatusInactive, err)
	assert.Nil(t, user)

	// nothing to change
	user, err = usecase.UpdateUser(repo, []string{"0", "test0", "TestUser0"})
	require.Equal(t, domain.ErrNothingToChange, err)
	assert.Nil(t, user)

	// success change
	user, err = usecase.UpdateUser(repo, []string{"2", "test22", "TestUser22"})
	require.NoError(t, err)
	assert.NotNil(t, user)

	assert.Equal(t, domain.User{Id: 2, Login: "test22", Name: "TestUser22", Status: true}, user[0])
}

func TestCurrentUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no users
	user, err := usecase.CurrentUser(repo)
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, user)

	// empty active file
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	user, err = usecase.CurrentUser(repo)
	require.NoError(t, err)
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}, user[0])

	// not empty active file
	user, err = usecase.CurrentUser(repo)
	require.NoError(t, err)
	assert.Equal(t, domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}, user[0])
}

func TestSwitchUser(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no args
	user, err := usecase.SwitchUser(repo, []string{})
	require.Equal(t, domain.ErrMissingActiveUserID, err)
	assert.Nil(t, user)

	// many args
	user, err = usecase.SwitchUser(repo, []string{"test0", "TestUser0"})
	require.Equal(t, domain.ErrTooManyArguments, err)
	assert.Nil(t, user)

	// one arg, not int
	user, err = usecase.SwitchUser(repo, []string{"two"})
	require.Error(t, err)
	assert.Nil(t, user)

	// one arg, int, no user exist
	user, err = usecase.SwitchUser(repo, []string{"2"})
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, user)

	// one arg, int
	_, err = usecase.AddUser(repo, []string{"test0", "TestUser0"})
	_, err = usecase.AddUser(repo, []string{"test1", "TestUser1"})
	_, err = usecase.AddUser(repo, []string{"test2", "TestUser2"})

	user, err = usecase.SwitchUser(repo, []string{"2"})
	require.NoError(t, err)
	assert.Equal(t, domain.User{Id: 2, Login: "test2", Name: "TestUser2", Status: true}, user[0])
}

func TestInfoUser(t *testing.T) {

}

func TestInfoAllActiveUsers(t *testing.T) {

}
