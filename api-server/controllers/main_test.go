package controllers_test

import (
	"testing"

	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers/testdata"
)

// テストに使用するコントローラ構造体を用意
var aCon *controllers.ArticleController

func TestMain(m *testing.M) {

	ser := testdata.NewServiceMock()
	aCon = controllers.NewArticleController(ser)

	// テスト実行
	m.Run()
}
