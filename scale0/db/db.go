package db

import (
	"context"

	"github.com/gerifield/service-scaling/scale0/model"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) Save(ctx context.Context, id string, content string) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO messages (id, content) VALUES (?, ?)", id, content)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetAll(ctx context.Context) ([]model.Message, error) {
	var resp []model.Message

	return resp, d.db.Select(&resp, "SELECT id, content, created_at FROM messages")
}
