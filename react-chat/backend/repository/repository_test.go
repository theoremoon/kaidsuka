package repository

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/theoremoon/kaidsuka/react-chat/backend/model"
)

const dbname = "database.db"

func setupRepository() Repository {
	repo, err := New(dbname)
	if err != nil {
		panic(err)
	}

	if err := repo.Setup(); err != nil {
		panic(err)
	}
	return repo
}
func teardownRepository(repo Repository) {
	repo.Close()
	os.RemoveAll(dbname)
}

func TestRepository_RegisterUser(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, "user1")

	_, err = repo.RegisterUser("user1")
	assert.Error(t, err, "")
}

func TestRepository_GetUser(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)

	user2, err := repo.GetUser("user1")
	assert.NoError(t, err)
	assert.Equal(t, user.Username, user2.Username)

	_, err = repo.GetUser("user000")
	assert.Error(t, err, "")
}

func TestRepository_PostMessage(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)
	_, err = repo.PostMessage(user.ID, "hello world")
	assert.NoError(t, err)
	_, err = repo.PostMessage(user.ID, "hello world")
	assert.NoError(t, err)
}

func TestRepository_GetMessage(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)
	msg, err := repo.PostMessage(user.ID, "hello world")
	assert.NoError(t, err)

	msg2, err := repo.GetMessage(msg.ID)
	assert.NoError(t, err)
	assert.Equal(t, msg2.Text, "hello world")
}

func TestRepository_UpdateMessage(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)
	msg, err := repo.PostMessage(user.ID, "hello world")
	assert.NoError(t, err)

	_, err = repo.UpdateMessage(user.ID, msg.ID, "goodbye world")
	assert.NoError(t, err)

	msg2, err := repo.GetMessage(msg.ID)
	assert.NoError(t, err)
	assert.Equal(t, msg2.Text, "goodbye world")
	assert.Equal(t, msg2.Edited, true)

	_, err = repo.UpdateMessage(0, msg.ID, "unmodifiable")
	assert.NoError(t, err)

	msg3, err := repo.GetMessage(msg.ID)
	assert.NoError(t, err)
	assert.Equal(t, msg3.Text, "goodbye world")
	assert.Equal(t, msg3.Edited, true)
}

func TestRepository_DeleteMessage(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)
	msg, err := repo.PostMessage(user.ID, "hello world")
	assert.NoError(t, err)

	err = repo.DeleteMessage(msg.ID)
	assert.NoError(t, err)

	_, err = repo.GetMessage(msg.ID)
	assert.Error(t, err)
}

func TestRepository_ListMessages(t *testing.T) {
	repo := setupRepository()
	defer teardownRepository(repo)

	user, err := repo.RegisterUser("user1")
	assert.NoError(t, err)

	const N = 20
	msgs := make([]*model.Message, N)
	for i := 0; i < N; i++ {
		msg, err := repo.PostMessage(user.ID, fmt.Sprintf("message %d", i+1))
		assert.NoError(t, err)
		msgs[i] = msg
		time.Sleep(1 * time.Second)
	}

	list, err := repo.ListMessages(math.MaxUint32, N/2)
	assert.NoError(t, err)
	for i := 0; i < N/2; i++ {
		assert.Equal(t, list[i].Text, msgs[N-i-1].Text)
	}

	list, err = repo.ListMessages(list[N/2-1].PostedAt, N/2)
	assert.NoError(t, err)
	for i := 0; i < N/2; i++ {
		assert.Equal(t, list[i].Text, msgs[N-i-1-(N/2)].Text)
	}
}
