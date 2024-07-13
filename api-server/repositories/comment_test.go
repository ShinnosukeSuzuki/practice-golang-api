package repositories_test

import (
	"testing"

	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/ShinnosukeSuzuki/practice-golang-api/repositories"

	_ "github.com/go-sql-driver/mysql"
)

// SelectCommentList関数のテスト
func TestSelectCommentList(t *testing.T) {
	articleID := 1
	got, err := repositories.SelectCommentList(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	for _, comment := range got {
		if comment.ArticleID != articleID {
			t.Errorf("want comment of articleID %d but got ID %d\n", articleID, comment.ArticleID)
		}
	}
}

// InsertComment関数のテスト
func TestInsertComment(t *testing.T) {
	comment := models.Comment{
		ArticleID: 1,
		Message:   "CommentInsertTest",
	}

	expectedCommentID := 3
	newComment, err := repositories.InserComment(testDB, comment)
	if err != nil {
		t.Fatal(err)
	}
	if newComment.CommentID != expectedCommentID {
		t.Errorf("new comment id is expected %d but got %d\n", expectedCommentID, newComment.CommentID)
	}

	// テストデータをCleanUp
	t.Cleanup(func() {
		const sqlDeleteComment = `
			DELETE FROM comments
			WHERE message = ?;
		`
		testDB.Exec(sqlDeleteComment, comment.Message)
	})
}
