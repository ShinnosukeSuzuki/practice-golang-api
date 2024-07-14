package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
	"github.com/ShinnosukeSuzuki/practice-golang-api/services"
	"github.com/gorilla/mux"
)

// コントローラ構造体を定義
type MyAppController struct {
	// フィールドにサービス構造体を持つ
	service *services.MyAppService
}

// コントローラのコンストラクタ
func NewMyAppController(s *services.MyAppService) *MyAppController {
	return &MyAppController{service: s}
}

func (c *MyAppController) HelloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, world!\n")
}

func (c *MyAppController) PostArticleHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article

	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
	}

	article, err := c.service.PostArticleService(reqArticle)
	if err != nil {
		http.Error(w, "fail to post article\n", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c *MyAppController) ArticleListHandler(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()

	var page int
	if p, ok := queryMap["page"]; ok && len(p) > 0 {
		var err error
		page, err = strconv.Atoi(p[0])
		if err != nil {
			http.Error(w, "Invalid query parameter", http.StatusBadRequest)
			return
		}
	} else {
		page = 1
	}

	articleList, err := c.service.GetArticleListService(page)
	if err != nil {
		http.Error(w, "fail to get article list\n", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(articleList)
}

func (c *MyAppController) ArticleDetailHandler(w http.ResponseWriter, r *http.Request) {
	articleID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := c.service.GetArticleService(articleID)
	if err != nil {
		http.Error(w, "fail to get article detail\n", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (c *MyAppController) PostNiceHandler(w http.ResponseWriter, r *http.Request) {
	var reqArticle models.Article
	if err := json.NewDecoder(r.Body).Decode(&reqArticle); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
	}

	article, err := c.service.PostNiceService(reqArticle)
	if err != nil {
		http.Error(w, "fail to post nice\n", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(article)
}

func (c *MyAppController) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reqComment); err != nil {
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		http.Error(w, "fail to post comment\n", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comment)
}
