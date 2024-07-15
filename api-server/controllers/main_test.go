package controllers_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/services"

	_ "github.com/go-sql-driver/mysql"
)

// テストに使用するコントローラ構造体を用意
var aCon *controllers.ArticleController

func TestMain(m *testing.M) {
	// DB接続
	dbUser := "docker"
	dbPassword := "docker"
	dbDatabase := "sampledb"
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("DB setup fail")
		os.Exit(1)
	}
	ser := services.NewMyAppService(db)
	aCon = controllers.NewArticleController(ser)

	// テスト実行
	m.Run()
}
