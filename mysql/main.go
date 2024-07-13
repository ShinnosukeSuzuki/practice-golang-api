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

	articleID := 100
	const sqlStr = `
		SELECT *
		FROM articles
		where article_id = ?;
	`

	row := db.QueryRow(sqlStr, articleID)
	if err := row.Err(); err != nil {
		// データ取得件数が0件の場合はデータ読み出し処理に移らずに終了
		fmt.Println(err)
		return
	}

	var article models.Article
	var createdTimme sql.NullTime

	// Scanメソッドで取得したデータを構造体に格納
	err = row.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTimme)
	if err != nil {
		fmt.Println(err)
		return
	}

	// createdTimmeがnullではない場合は、article.CreatedAtに格納
	if createdTimme.Valid {
		article.CreatedAt = createdTimme.Time
	}

	fmt.Printf("%+v\n", article)
}
