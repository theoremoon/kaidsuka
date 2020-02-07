package model

type User struct {
	ID       uint32 `db:"id"`
	Username string `db:"username"`
}

type Message struct {
	ID       uint32 `db:"id"`
	UserID   uint32 `db:"user_id"`
	Text     string `db:"text"`
	Edited   bool   `db:"edited"`
	PostedAt uint32 `db:"posted_at"`
}
