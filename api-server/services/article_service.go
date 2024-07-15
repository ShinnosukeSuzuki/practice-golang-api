package services

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/ShinnosukeSuzuki/practice-golang-api/apperrors"
	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/ShinnosukeSuzuki/practice-golang-api/repositories"
)

// ハンドラ ArticleDetailHandler 用のサービス関数
func (s *MyAppService) GetArticleService(articleID int) (models.Article, error) {
	// 並列処理でarticleとcommentListを取得

	var article models.Article
	var commentList []models.Comment
	var articleGetErr, commentGetErr error

	// それぞれをMutexで保護
	var amu sync.Mutex
	var cmu sync.Mutex

	// メインゴールーチンの終了を待つ
	var wg sync.WaitGroup
	wg.Add(2)

	// articleIDに一致する記事の詳細情報を取得
	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		amu.Lock()
		article, articleGetErr = repositories.SelectArticleDetail(db, articleID)
		amu.Unlock()
	}(s.db, articleID)

	// articleIDに一致する記事のコメント情報を取得
	go func(db *sql.DB, articleID int) {
		defer wg.Done()
		cmu.Lock()
		commentList, commentGetErr = repositories.SelectCommentList(db, articleID)
		cmu.Unlock()
	}(s.db, articleID)

	// 両方のゴールーチンが終了するまで待つ
	wg.Wait()

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
