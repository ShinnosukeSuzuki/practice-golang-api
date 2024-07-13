package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
)

const (
	articleNumPerPage = 5
)

// 新規投稿をデータベースにinsertする関数
func InsertArticle(db *sql.DB, article models.Article) (models.Article, error) {
	const sqlInsertArticle = `
		INSERT INTO articles (title, contents, username, nice, created_at)
		VALUES (?, ?, ?, 0, NOW());
	`

	var newArticle models.Article
	newArticle.Title = article.Title
	newArticle.Contents = article.Contents
	newArticle.UserName = article.UserName
	result, err := db.Exec(sqlInsertArticle, article.Title, article.Contents, article.UserName)
	if err != nil {
		return models.Article{}, err
	}
	id, _ := result.LastInsertId()
	newArticle.ID = int(id)

	return newArticle, nil
}

// 変数pageで指定されたページに表示する投稿一覧を取得する関数
func SelectArticleList(db *sql.DB, page int) ([]models.Article, error) {
	const sqlSelectArticleList = `
		SELECT article_id, title, contents, username, nice, created_at
		FROM articles
		LIMIT ? OFFSET ?;
	`

	offset := (page - 1) * articleNumPerPage
	rows, err := db.Query(sqlSelectArticleList, articleNumPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 取得したデータを格納するスライス
	articleArray := make([]models.Article, 0)
	for rows.Next() {
		var article models.Article
		// create_atがnullの対策
		var createdTime sql.NullTime

		rows.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTime)

		if createdTime.Valid {
			article.CreatedAt = createdTime.Time
		}
		articleArray = append(articleArray, article)
	}

	return articleArray, nil
}

// 投稿IDを指定して、記事データを取得する関数
func SelectArticleDetail(db *sql.DB, articleID int) (models.Article, error) {
	const sqlSelectArticleDetail = `
		SELECT *
		FROM articles
		WHERE article_id = ?;
	`

	row := db.QueryRow(sqlSelectArticleDetail, articleID)
	if err := row.Err(); err != nil {
		return models.Article{}, err
	}

	var article models.Article
	var createdTime sql.NullTime

	err := row.Scan(&article.ID, &article.Title, &article.Contents, &article.UserName, &article.NiceNum, &createdTime)
	if err != nil {
		return models.Article{}, err
	}

	if createdTime.Valid {
		article.CreatedAt = createdTime.Time
	}

	return article, nil
}

// いいねの数を1増やす関数
func UpdateNiceNum(db *sql.DB, articleID int) error {
	// 3-7でやったようにトランザクションを使い、いいね数を取得→更新→コミットする
	// トランザクションの開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 現在のいいね数を取得するクエリを実行する
	const sqlGetNice = `
		SELECT nice
		FROM articles
		WHERE article_id = ?;
	`

	row := tx.QueryRow(sqlGetNice, articleID)
	if err := row.Err(); err != nil {
		tx.Rollback()
		return err
	}

	var nicenum int
	err = row.Scan(&nicenum)
	if err != nil {
		tx.Rollback()
		return err
	}

	// いいね数を1増やすクエリを実行する
	const sqlUpdateNice = `
		UPDATE articles
		SET nice = ?
		WHERE article_id = ?;
	`

	_, err = tx.Exec(sqlUpdateNice, nicenum+1, articleID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// トランザクションのコミット
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
