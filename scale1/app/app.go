package app

import (
	"context"
	"log"
	"time"

	"github.com/gerifield/service-scaling/scale1/model"
	"github.com/google/uuid"
)

const getAllCacheKey = "cacheKey:getAll"

var getAllCacheTTL = time.Second * 10

type messageRepo interface {
	Save(ctx context.Context, id string, content string) error
	GetAll(ctx context.Context) ([]model.Message, error)
}

type cache interface {
	GetMessages(ctx context.Context, key string) ([]model.Message, error)
	SaveMessages(ctx context.Context, key string, value []model.Message, ttl time.Duration) error
	Invalidate(ctx context.Context, key string) error
}

type Logic struct {
	messageRepo messageRepo
	cache       cache
}

func New(messageRepo messageRepo, cache cache) *Logic {
	return &Logic{
		messageRepo: messageRepo,
		cache:       cache,
	}
}

func (l *Logic) Save(ctx context.Context, content string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	err = l.messageRepo.Save(ctx, id.String(), content)
	if err == nil {
		log.Println("invalidating the cache")
		_ = l.cache.Invalidate(ctx, getAllCacheKey)
	}

	return id.String(), err
}

func (l *Logic) GetAll(ctx context.Context) ([]model.Message, error) {
	res, err := l.cache.GetMessages(ctx, getAllCacheKey)
	if err == nil {
		log.Println("return cached result")
		return res, nil
	}

	messages, err := l.messageRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("save values into cache")
	err = l.cache.SaveMessages(ctx, getAllCacheKey, messages, getAllCacheTTL)
	if err != nil {
		log.Println("cache save failed", err)
	}

	log.Println("return db result")
	return messages, nil
}
