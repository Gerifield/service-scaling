package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	queueReadTimeout = 5 * time.Second
)

type Queue struct {
	queueName   string
	redisClient *redis.Client
}

func New(redisClient *redis.Client, queueName string) *Queue {
	return &Queue{
		redisClient: redisClient,
		queueName:   queueName,
	}
}

func (q *Queue) Put(ctx context.Context, val string) error {
	return q.redisClient.RPush(ctx, q.queueName, val).Err()
}

func (q *Queue) Get(ctx context.Context) (string, error) {
	res := q.redisClient.BLPop(ctx, queueReadTimeout, q.queueName)
	if res.Err() != nil {
		return "", res.Err()
	}

	r := res.Val()
	if len(r) < 2 {
		return "", fmt.Errorf("not enough param %v", r)
	}

	return r[1], nil
}
