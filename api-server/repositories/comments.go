package repositories

import (
	"database/sql"

	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
)

// 新規コメントをデータベースにinsertする関数
func InserComment(db *sql.DB, comment models.Comment) (models.Comment, error) {
	const sqlInsertComment = `
		INSERT INTO comments (article_id, message, created_at)
		VALUES (?, ?, NOW());
	`

	var newComment models.Comment
	newComment.ArticleID = comment.ArticleID
	newComment.Message = comment.Message

	result, err := db.Exec(sqlInsertComment, comment.ArticleID, comment.Message)
	if err != nil {
		return models.Comment{}, err
	}
	id, _ := result.LastInsertId()
	newComment.CommentID = int(id)

	return newComment, nil
}

// 指定IDの記事についたコメント一覧を取得する関数
func SelectCommentList(db *sql.DB, articleID int) ([]models.Comment, error) {
	// article_idが一致するコメントを取得するクエリ
	const sqlSelectCommentList = `
		SELECT *
		FROM comments
		WHERE article_id = ?;
	`

	rows, err := db.Query(sqlSelectCommentList, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 取得したデータを格納するスライス
	commentArray := make([]models.Comment, 0)
	for rows.Next() {
		var comment models.Comment
		var createdTime sql.NullTime

		rows.Scan(&comment.CommentID, &comment.ArticleID, &comment.Message, &createdTime)

		if createdTime.Valid {
			comment.CreatedAt = createdTime.Time
		}

		commentArray = append(commentArray, comment)
	}

	return commentArray, nil
}
