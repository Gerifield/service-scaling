package main

import (
	"log"
	"net/http"

	"github.com/gerifield/service-scaling/scale4/app"
	"github.com/gerifield/service-scaling/scale4/cache"
	"github.com/gerifield/service-scaling/scale4/db"
	"github.com/gerifield/service-scaling/scale4/queue"
	"github.com/gerifield/service-scaling/scale4/server"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/namsral/flag"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP Listen address")
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

	redisCache := cache.NewRedis(rdb)
	msgSaveQueue := queue.New(rdb, "queue:message-save")
	repo := db.New(wDBConn, rDBConn)
	appLogic := app.New(repo, redisCache, msgSaveQueue)
	s := server.New(appLogic)

	log.Println("Started", *addr)
	http.ListenAndServe(*addr, s.Routes())
}
