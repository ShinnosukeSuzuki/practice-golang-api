package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ShinnosukeSuzuki/practice-golang-api/apperrors"
	"github.com/ShinnosukeSuzuki/practice-golang-api/common"
	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers/services"
	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/gorilla/mux"
)

// Article用のコントローラ構造体
type ArticleController struct {
	service services.ArticleServicer
}

// コンストラクタ関数
func NewArticleController(s services.ArticleServicer) *ArticleController {
	return &ArticleController{service: s}
}

// ハンドラメソッドを定義
// GET /hello
func (c *ArticleController) HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

// POST /article
func (c *ArticleController) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article

	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	authedUsername := common.GetUserName(r.Context())
	if reqArticle.UserName != authedUsername {
		err := apperrors.NotMatchUser.Wrap(errors.New("does not match reqBody user and idtoken user"), "invalid user")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	article, err := c.service.PostArticleService(reqArticle)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// GET /article/list
func (c *ArticleController) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()

	var page int
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			err = apperrors.BadParam.Wrap(err, "queryparam must be number")
			apperrors.ErrorHandler(w, r, err)
			return
		}
	} else {
		page = 1
	}

	articleList, err := c.service.GetArticleListService(page)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(articleList)
}

// GET /article/{id}
func (c *ArticleController) ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		err = apperrors.BadParam.Wrap(err, "pathparam must be number")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	article, err := c.service.GetArticleService(articleID)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(article)
}

// POST /article/nice
func (c *ArticleController) PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
		return
	}

	article, err := c.service.PostNiceService(reqArticle)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(article)
}
