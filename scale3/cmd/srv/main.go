package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gerifield/service-scaling/scale2/app"
	"github.com/gerifield/service-scaling/scale2/cache"
	"github.com/gerifield/service-scaling/scale2/db"
	"github.com/gerifield/service-scaling/scale2/server"
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
	repo := db.New(wDBConn, rDBConn)
	appLogic := app.New(repo, redisCache)
	s := server.New(appLogic)

	log.Println("Started", *addr)
	http.ListenAndServe(*addr, s.Routes())
}
