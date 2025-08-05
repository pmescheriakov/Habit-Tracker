package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pmescheriakov/Habit-Tracker/internal/domain"
	"github.com/pmescheriakov/Habit-Tracker/internal/infra/db"
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

func TestGetAll(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// empty file
	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 0, len(users))

	// one user
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	users, err = repo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, user0.Login, users[0].Login)
	assert.Equal(t, user0.Name, users[0].Name)

	// some users
	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: false}
	err = repo.Save(user1)
	require.NoError(t, err)

	users, err = repo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 2, len(users))

	logins := []string{users[0].Login, users[1].Login}

	expectedLogins := map[string]bool{
		"test0": true,
		"test1": true,
	}

	for _, login := range logins {
		assert.True(t, expectedLogins[login])
	}
}

func TestFindLogName(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// empty file
	pUser, err := repo.FindLogName("test0", "TestUser0")
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, pUser)

	// no empty file
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	// success find
	pUser, err = repo.FindLogName("test1", "TestUser1")
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, user1.Login, pUser.Login)
	assert.Equal(t, user1.Name, pUser.Name)

	// unsuccess find
	pUser, err = repo.FindLogName("test3", "TestUser3")
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, pUser)
}

func TestFindId(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// empty file
	pUser, err := repo.FindId(0)
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, pUser)

	// no empty file
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	// success find
	pUser, err = repo.FindId(1)
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, user1.Login, pUser.Login)
	assert.Equal(t, user1.Name, pUser.Name)

	// unsuccess find
	pUser, err = repo.FindId(2)
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, pUser)
}

func TestSave(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}

	err := repo.Save(user0)
	require.NoError(t, err)

	users, err := repo.GetAll()
	require.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, user0.Login, users[0].Login)
	assert.Equal(t, user0.Name, users[0].Name)
}

func TestUpdate(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// empty file
	usersOld, err := repo.GetAll()
	require.NoError(t, err)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Update(user0)
	require.Equal(t, domain.ErrUserNotFound, err)

	usersNew, err := repo.GetAll()
	require.NoError(t, err)
	require.Equal(t, usersOld, usersNew)

	// no empty file
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	// success update
	user1Upd := domain.User{Id: 1, Login: "test111", Name: "TestUser11111", Status: true}
	err = repo.Update(user1Upd)
	require.NoError(t, err)

	pUser, err := repo.FindId(user1Upd.Id)
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, user1Upd.Login, pUser.Login)
	assert.Equal(t, user1Upd.Name, pUser.Name)

	// unsuccess update
	user3Upd := domain.User{Id: 3, Login: "test33", Name: "TestUser3333", Status: true}
	err = repo.Update(user3Upd)
	require.Equal(t, domain.ErrUserNotFound, err)

	pUser, err = repo.FindId(user3Upd.Id)
	require.Equal(t, domain.ErrUserNotFound, err)
	assert.Nil(t, pUser)
}

func TestActivate(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no user
	err := repo.Activate(1)
	require.Equal(t, domain.ErrUserNotFound, err)

	// user exists
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: false}
	err = repo.Save(user1)
	require.NoError(t, err)

	err = repo.Activate(user1.Id)
	require.NoError(t, err)

	userAfterDeactivate, err := repo.FindId(user1.Id)
	require.NoError(t, err)
	assert.Equal(t, !user1.Status, userAfterDeactivate.Status)
}

func TestDeactivate(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no user
	err := repo.Deactivate(1)
	require.Equal(t, domain.ErrUserNotFound, err)

	// user exists
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	err = repo.Deactivate(user1.Id)
	require.NoError(t, err)

	userAfterDeactivate, err := repo.FindId(user1.Id)
	require.NoError(t, err)
	assert.Equal(t, !user1.Status, userAfterDeactivate.Status)
}

func TestSetActive(t *testing.T) {
	repo, _, _ := setupTestRepo(t)

	// no user
	err := repo.SetActive(1)
	require.Equal(t, domain.ErrUserNotFound, err)

	// user exists
	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	err = repo.SetActive(1)
	require.NoError(t, err)

	userActId, err := repo.GetActive()
	require.NoError(t, err)
	assert.Equal(t, *userActId, user1)
}

func TestGetActive(t *testing.T) {
	// no users exists
	repo, _, _ := setupTestRepo(t)

	// empty active file
	pUser, err := repo.GetActive()
	require.Equal(t, domain.ErrNoUsers, err)
	assert.Nil(t, pUser)

	// users exists
	repo, _, _ = setupTestRepo(t)

	user0 := domain.User{Id: 0, Login: "test0", Name: "TestUser0", Status: true}
	err = repo.Save(user0)
	require.NoError(t, err)

	user1 := domain.User{Id: 1, Login: "test1", Name: "TestUser1", Status: true}
	err = repo.Save(user1)
	require.NoError(t, err)

	// empty active file
	pUser, err = repo.GetActive()
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, 0, pUser.Id)

	// init active file
	pUser, err = repo.GetActive()
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, 0, pUser.Id)

	//
	err = repo.SetActive(1)
	require.NoError(t, err)

	pUser, err = repo.GetActive()
	require.NoError(t, err)
	require.NotNil(t, pUser)
	assert.Equal(t, 1, pUser.Id)
}
