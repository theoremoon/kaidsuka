package repository

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rakyll/statik/fs"
	"github.com/theoremoon/kaidsuka/react-chat/backend/model"
	_ "github.com/theoremoon/kaidsuka/react-chat/backend/statik"
)

type Repository interface {
	Setup() error

	GetUser(username string) (*model.User, error)
	GetUserByID(id uint32) (*model.User, error)
	RegisterUser(username string) (*model.User, error)

	GetMessage(id uint32) (*model.Message, error)
	ListMessages(olderThan uint32, limit uint32) ([]*model.Message, error)
	PostMessage(userID uint32, text string) (*model.Message, error)
	UpdateMessage(userID, messageID uint32, text string) (*model.Message, error)
	DeleteMessage(id uint32) error

	Close() error
}

type repository struct {
	db *sqlx.DB
}

func New(database string) (Repository, error) {
	db, err := sqlx.Open("sqlite3", database)
	if err != nil {
		return nil, err
	}
	return &repository{
		db: db,
	}, nil
}

func (r *repository) GetUser(username string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user,
		`SELECT id, username
		FROM users
		WHERE username = ?`,
		username)
	return &user, err
}
func (r *repository) GetUserByID(id uint32) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user,
		`SELECT id, username
		FROM users
		WHERE id = ?`,
		id)
	return &user, err
}

func (r *repository) RegisterUser(username string) (*model.User, error) {
	id := r.generateID()
	_, err := r.db.Exec(
		`INSERT INTO users(id, username)
		VALUES (?, ?)`,
		id, username)
	if err != nil {
		return nil, err
	}
	return r.GetUser(username)
}

func (r *repository) GetMessage(id uint32) (*model.Message, error) {
	var message model.Message
	err := r.db.Get(&message,
		`SELECT id, user_id, text, edited, posted_at
		FROM messages
		WHERE id = ?`,
		id)
	return &message, err
}

func (r *repository) ListMessages(olderThan uint32, limit uint32) ([]*model.Message, error) {
	messages := make([]*model.Message, 0, limit)
	err := r.db.Select(
		&messages,
		`SELECT id, user_id, text, edited, posted_at
		FROM messages
		WHERE posted_at < ?
		ORDER BY posted_at DESC
		LIMIT ?
		`, olderThan, limit)
	return messages, err
}

func (r *repository) PostMessage(userID uint32, text string) (*model.Message, error) {
	id := r.generateID()
	_, err := r.db.Exec(
		`INSERT INTO messages(id, user_id, text, edited, posted_at)
		VALUES (?, ?, ?, 0, ?)`,
		id, userID, text, time.Now().Unix())
	if err != nil {
		return nil, err
	}
	return r.GetMessage(id)
}
func (r *repository) UpdateMessage(userID, messageID uint32, text string) (*model.Message, error) {
	_, err := r.db.Exec(
		`UPDATE messages
		SET text = ?, edited = 1
		WHERE id = ? AND user_id = ?`,
		text, messageID, userID)
	if err != nil {
		return nil, err
	}
	return r.GetMessage(messageID)
}

func (r *repository) DeleteMessage(id uint32) error {
	_, err := r.db.Exec(
		`DELETE FROM messages
		WHERE id = ?`, id)
	return err
}

func (r *repository) Setup() error {
	hfs, err := fs.New()
	if err != nil {
		return err
	}
	schemaFile, err := hfs.Open("/schema.sql")
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	schema, err := ioutil.ReadAll(schemaFile)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(string(schema))

	return err
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) generateID() uint32 {
	return rand.Uint32()
}
