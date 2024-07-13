package repositories_test

import (
	"testing"

	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/ShinnosukeSuzuki/practice-golang-api/repositories"
	"github.com/ShinnosukeSuzuki/practice-golang-api/repositories/testdata"

	_ "github.com/go-sql-driver/mysql"
)

// SelectArticleList関数のテスト
func TestSelectArticleList(t *testing.T) {
	expectedNum := len(testdata.ArticleTrstData)
	got, err := repositories.SelectArticleList(testDB, 1)
	if err != nil {
		t.Fatal(err)
	}

	if num := len(got); num != expectedNum {
		t.Errorf("want %d but got %d articles\n", expectedNum, num)
	}
}

// SelectArticleDetail関数のテスト
func TestSelectArticleDetail(t *testing.T) {
	// テストデータの投稿を作成
	tests := []struct {
		testTitle string
		expected  models.Article
	}{
		{
			testTitle: "subtest1",
			expected:  testdata.ArticleTrstData[0],
		}, {
			testTitle: "subtest2",
			expected:  testdata.ArticleTrstData[1],
		},
	}

	for _, test := range tests {
		// 各サブテストを実行
		t.Run(test.testTitle, func(t *testing.T) {
			got, err := repositories.SelectArticleDetail(testDB, test.expected.ID)
			if err != nil {
				t.Fatal(err)
			}

			if got.ID != test.expected.ID {
				t.Errorf("ID: get %d but want %d\n", got.ID, test.expected.ID)
			}
			if got.Title != test.expected.Title {
				t.Errorf("Title: get %s but want %s\n", got.Title, test.expected.Title)
			}
			if got.Contents != test.expected.Contents {
				t.Errorf("Content: get %s but want %s\n", got.Contents, test.expected.Contents)
			}
			if got.UserName != test.expected.UserName {
				t.Errorf("UserName: get %s but want %s\n", got.UserName, test.expected.UserName)
			}
			if got.NiceNum != test.expected.NiceNum {
				t.Errorf("NiceNum: get %d but want %d\n", got.NiceNum, test.expected.NiceNum)
			}
		})
	}
}

// InsertArticle 関数のテスト
func TestInsertArticle(t *testing.T) {
	article := models.Article{
		Title:    "insertTest",
		Contents: "testest",
		UserName: "saki",
	}

	expectedArticleNum := 3
	newArticle, err := repositories.InsertArticle(testDB, article)
	if err != nil {
		t.Error(err)
	}
	if newArticle.ID != expectedArticleNum {
		t.Errorf("new article id is expected %d but got %d\n", expectedArticleNum, newArticle.ID)
	}

	// 他のテストに影響しないようにcleanupする
	t.Cleanup(func() {
		// Insertしたデータを削除するクエリ
		const sqlDeleteArticle = `
			DELETE FROM articles
			WHERE title = ? and contents = ? and username = ?;
		`
		testDB.Exec(sqlDeleteArticle, article.Title, article.Contents, article.UserName)
	})
}

// UpdateNiceNum 関数のテスト
func TestUpdateNiceNum(t *testing.T) {
	articleID := 1
	// いいね数を1増やす
	err := repositories.UpdateNiceNum(testDB, articleID)
	if err != nil {
		t.Fatal(err)
	}

	// articleID=1の記事の詳細を取得
	got, _ := repositories.SelectArticleDetail(testDB, articleID)

	// 更新後のデータと元のテストデータのいいね数の差が1でない場合はエラー
	if got.NiceNum-testdata.ArticleTrstData[articleID-1].NiceNum != 1 {
		t.Errorf("fail to update nice num: expected %d but got %d\\n", testdata.ArticleTrstData[articleID-1].NiceNum+1, got.NiceNum)
	}
}
