package api

import (
	"database/sql"
	"net/http"

	"github.com/ShinnosukeSuzuki/practice-golang-api/api/middlewares"
	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers"
	"github.com/ShinnosukeSuzuki/practice-golang-api/services"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	ser := services.NewMyAppService(db)
	aCon := controllers.NewArticleController(ser)
	cCon := controllers.NewCommentController(ser)

	// ルータの初期化
	r := mux.NewRouter()

	// ルーティングの設定
	r.HandleFunc("/hello", aCon.HelloHandler).Methods(http.MethodGet)
	r.HandleFunc("/article", aCon.PostArticleHandler).Methods(http.MethodPost)
	r.HandleFunc("/article/list", aCon.ArticleListHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/{id:[0-9]+}", aCon.ArticleDetailHandler).Methods(http.MethodGet)
	r.HandleFunc("/article/nice", aCon.PostNiceHandler).Methods(http.MethodPost)
	r.HandleFunc("/comment", cCon.PostCommentHandler).Methods(http.MethodPost)

	r.Use(middlewares.LoggingMiddleware)
	r.Use(middlewares.AuthMiddleware)

	return r
}
