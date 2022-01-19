package model

import "time"

type Message struct {
	ID        string    `db:"id" json:"id"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
