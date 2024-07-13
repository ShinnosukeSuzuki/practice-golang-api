package main

import (
	"database/sql"
	"dbsample/models"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := "docker"
	dbPassword := "docker"
	dbDatabase := "sampledb"
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	article := models.Article{
		Title:    "insert test",
		Contents: "Can I insert data correctly?",
		UserName: "hogehoge",
	}
	const sqlStr = `
		INSERT INTO articles (title, contents, username, nice, created_at)
		VALUES (?, ?, ?, 0, NOW());
	`

	result, err := db.Exec(sqlStr, article.Title, article.Contents, article.UserName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
}
