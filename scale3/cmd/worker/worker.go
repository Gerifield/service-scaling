package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"

	"github.com/gerifield/service-scaling/scale3/cache"
	"github.com/gerifield/service-scaling/scale3/db"
	"github.com/gerifield/service-scaling/scale3/model"
	"github.com/gerifield/service-scaling/scale3/queue"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	wSQLConf := flag.String("writeSQLConf", "api:api@tcp(db:3306)/api?parseTime=true", "Writer MySQL connection string")
	rSQLConf := flag.String("readSQLConf", "api:api@tcp(db:3306)/api?parseTime=true", "Reader MySQL connection string")
	redisAddr := flag.String("redisAddr", "redis:6379", "Redis connection address")
	flag.Parse()

	wDBConn, err := sqlx.Connect("mysql", *wSQLConf)
	if err != nil {
		log.Fatalln(err)
	}

	rDBConn, err := sqlx.Connect("mysql", *rSQLConf)
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	msgSaveQueue := queue.New(rdb, "queue:message-save")
	repo := db.New(wDBConn, rDBConn)
	redisCache := cache.NewRedis(rdb)

	for {
		log.Println("Wait for message")
		it, err := msgSaveQueue.Get(context.Background())
		if err != nil {

			log.Println("Queue read failed", err)
			continue
		}

		var pl model.QueueMessage
		err = json.Unmarshal([]byte(it), &pl)
		if err != nil {
			log.Println("JSON read failed", err)
			continue
		}

		switch pl.MsgType {
		case model.MsgTypeSave:
			payload := pl.Payload

			err = repo.Save(context.Background(), payload.ID, payload.Content)
			if err != nil {
				log.Println("db save failed", err)
				continue
			}
			log.Println("Message saved", payload.ID)

			log.Println("invalidating the cache")
			_ = redisCache.Invalidate(context.Background(), "cacheKey:getAll")
		}
	}
}
