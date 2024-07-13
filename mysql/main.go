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

	const sqlStr = `
		SELECT *
		FROM articles;
	`
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	articleArrray := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		var createdTimme sql.NullTime
		err := rows.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTimme)

		if createdTimme.Valid {
			article.CreatedAt = createdTimme.Time
		}

		if err != nil {
			fmt.Println(err)
		} else {
			articleArrray = append(articleArrray, article)
		}
	}

	fmt.Printf("%+v\n", articleArrray)
}
