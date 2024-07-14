package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/routers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/services"

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

	// sql.DB型を元に、サーバー全体で使用するサービス構造体MyAppServiceを1つ生成
	ser := services.NewMyAppService(db)

	// MyAppService型を元に、コントローラー全体で使用するコントローラ構造体MyAppControllerを1つ生成
	con := controllers.NewMyAppController(ser)

	// コントローラ型MyAppControllerのハンドラメソッドとパスを紐付け
	r := routers.NewRouter(con)

	log.Println("server start at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
