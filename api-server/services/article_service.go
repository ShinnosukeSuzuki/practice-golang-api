package services

import (
	"database/sql"
	"errors"

	"github.com/ShinnosukeSuzuki/practice-golang-api/apperrors"
	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/ShinnosukeSuzuki/practice-golang-api/repositories"
)

// ハンドラ ArticleDetailHandler 用のサービス関数
func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	// 並行処理でarticleとcommentListを取得

	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	// Article型とerror型を同時に扱うための構造体
	type articleResult struct {
		article models.Article
		err     error
	}
	// articleResult型のチャネルを定義
	articleChan := make(chan articleResult)
	defer close(articleChan)

	// articleIDに一致する記事の詳細情報を取得
	// articleChanを通じて、SelectArticleDetail関数の結果を送信
	go func(ch chan<- articleResult, db *sql.DB, articleID int) {
		article, articleGetErr = repositories.SelectArticleDetail(db, articleID)
		ch <- articleResult{article: article, err: articleGetErr}
	}(articleChan, s.db, articleID)

	// Comment型とerror型を同時に扱うための構造体
	type commentResult struct {
		commentList []models.Comment
		err         error
	}
	// commentResult型のチャネルを定義
	commentChan := make(chan commentResult)
	defer close(commentChan)

	// articleIDに一致する記事のコメント情報を取得
	// commentChanを通じて、SelectCommentList関数の結果を送信
	go func(ch chan<- commentResult, db *sql.DB, articleID int) {
		commentList, commentGetErr = repositories.SelectCommentList(db, articleID)
		ch <- commentResult{commentList: commentList, err: commentGetErr}
	}(commentChan, s.db, articleID)

	for i := 0; i < 2; i++ {
		select {
		case ar := <-articleChan:
			article, articleGetErr = ar.article, ar.err
		case cr := <-commentChan:
			commentList, commentGetErr = cr.commentList, cr.err
		}
	}

	if articleGetErr != nil {
		if errors.Is(articleGetErr, sql.ErrNoRows) {
			err := apperrors.NAData.Wrap(articleGetErr, "no data")
			return models.Article{}, err
		}
		err := apperrors.GetDataFailed.Wrap(articleGetErr, "fail to get data")
		return models.Article{}, err
	}

	if commentGetErr != nil {
		err := apperrors.GetDataFailed.Wrap(commentGetErr, "fail to get data")
		return models.Article{}, err
	}

	// 取得したコメント情報を記事情報に追加
	article.CommentList = append(article.CommentList, commentList...)
	return article, nil

}

// PostArticleHandler で使うことを想定したサービス
// 引数の情報をもとに新しい記事を作り、結果を返却
func (s *MyAppService) PostArticleService(article models.Article) (models.Article, error) {
	newArticle, err := repositories.InsertArticle(s.db, article)
	if err != nil {
		err = apperrors.InsertDataFailed.Wrap(err, "Failed to record data")
		return models.Article{}, err
	}

	return newArticle, nil
}

// ArticleListHandler で使うことを想定したサービス
// 指定 page の記事一覧を返却
func (s *MyAppService) GetArticleListService(page int) ([]models.Article, error) {
	articleList, err := repositories.SelectArticleList(s.db, page)
	if err != nil {
		err = apperrors.GetDataFailed.Wrap(err, "Failed to get data")
		return nil, err
	}

	// articleListが空の場合
	if len(articleList) == 0 {
		err := apperrors.NAData.Wrap(ErrNoData, "no data")
		return nil, err
	}

	return articleList, nil
}

// PostNiceHandler で使うことを想定したサービス
// 指定 ID の記事のいいね数を+1 して、結果を返却
func (s *MyAppService) PostNiceService(article models.Article) (models.Article, error) {
	err := repositories.UpdateNiceNum(s.db, article.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = apperrors.NAData.Wrap(err, "does not exist target article")
			return models.Article{}, err
		}
		err = apperrors.UpdateDataFailed.Wrap(err, "Failed to update nice count")
		return models.Article{}, err
	}

	return models.Article{
		ID:        article.ID,
		Title:     article.Title,
		Contents:  article.Contents,
		UserName:  article.UserName,
		NiceNum:   article.NiceNum + 1,
		CreatedAt: article.CreatedAt,
	}, nil
}
