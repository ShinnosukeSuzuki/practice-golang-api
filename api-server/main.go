package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShinnosukeSuzuki/practice-golang-api/api"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbDatabase = os.Getenv("DB_NAME")
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

func main() {
	// DB接続
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("failed to connect database")
		return
	}

	// コントローラ型MyAppControllerのハンドラメソッドとパスを紐付け
	r := api.NewRouter(db)

	log.Println("server start at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
