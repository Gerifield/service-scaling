package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gerifield/service-scaling/scale1/app"
	"github.com/gerifield/service-scaling/scale1/cache"
	"github.com/gerifield/service-scaling/scale1/db"
	"github.com/gerifield/service-scaling/scale1/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP Listen address")
	sqlConf := flag.String("sqlconf", "api:api@tcp(db:3306)/api?parseTime=true", "MySQL connection string")
	redisAddr := flag.String("redisAddr", "redis:6379", "Redis connection address")
	flag.Parse()

	dbConn, err := sqlx.Connect("mysql", *sqlConf)
	if err != nil {
		log.Fatalln(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     *redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	redisCache := cache.NewRedis(rdb)
	repo := db.New(dbConn)
	appLogic := app.New(repo, redisCache)
	s := server.New(appLogic)

	log.Println("Started", *addr)
	http.ListenAndServe(*addr, s.Routes())
}
