package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/gerifield/service-scaling/scale0/app"
	"github.com/gerifield/service-scaling/scale0/db"
	"github.com/gerifield/service-scaling/scale0/server"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP Listen address")
	sqlConf := flag.String("sqlconf", "api:api@tcp(db:3306)/api?parseTime=true", "MySQL connection string")
	flag.Parse()

	dbConn, err := sqlx.Connect("mysql", *sqlConf)
	if err != nil {
		log.Fatalln(err)
	}

	repo := db.New(dbConn)
	appLogic := app.New(repo)
	s := server.New(appLogic)

	log.Println("Started", *addr)
	http.ListenAndServe(*addr, s.Routes())
}
