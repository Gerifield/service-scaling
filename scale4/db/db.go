package db

import (
	"context"

	"github.com/gerifield/service-scaling/scale4/model"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	writer *sqlx.DB
	reader *sqlx.DB
}

func New(writer, reader *sqlx.DB) *DB {
	return &DB{
		writer: writer,
		reader: reader,
	}
}

func (d *DB) Save(ctx context.Context, id string, content string) error {
	_, err := d.writer.ExecContext(ctx, "INSERT INTO messages (id, content) VALUES (?, ?)", id, content)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetAll(ctx context.Context) ([]model.Message, error) {
	resp := make([]model.Message, 0)

	return resp, d.reader.Select(&resp, "SELECT id, content, created_at FROM messages")
}
