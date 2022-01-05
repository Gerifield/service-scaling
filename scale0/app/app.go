package app

import (
	"context"

	"github.com/gerifield/service-scaling/scale0/model"
	"github.com/google/uuid"
)

type messageRepo interface {
	Save(ctx context.Context, id string, content string) error
	GetAll(ctx context.Context) ([]model.Message, error)
}

type Logic struct {
	messageRepo messageRepo
}

func New(messageRepo messageRepo) *Logic {
	return &Logic{
		messageRepo: messageRepo,
	}
}

func (l *Logic) Save(ctx context.Context, content string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), l.messageRepo.Save(ctx, id.String(), content)
}

func (l *Logic) GetAll(ctx context.Context) ([]model.Message, error) {
	return l.messageRepo.GetAll(ctx)
}
