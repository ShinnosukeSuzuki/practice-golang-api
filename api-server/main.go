package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/services"
	"github.com/gorilla/mux"

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
	r := mux.NewRouter()

	r.HandleFunc("/hello", con.HelloHandler).Methods(http.MethodGet)
	r.HandleFunc("/article", con.PostArticleHandler).Methods(http.MethodPost)
	r.HandleFunc("/article/list", con.ArticleListHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/{id:[0-9]+}", con.ArticleDetailHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/nice", con.PostNiceHandler).Methods(http.MethodPost)
	r.HandleFunc("/comment", con.PostCommentHandler).Methods(http.MethodPost)

	log.Println("server start at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
