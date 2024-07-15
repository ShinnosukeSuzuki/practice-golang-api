package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ShinnosukeSuzuki/practice-golang-api/apperrors"
	"github.com/ShinnosukeSuzuki/practice-golang-api/controllers/services"
	"github.com/ShinnosukeSuzuki/practice-golang-api/models"
)

// Comment用のコントローラ構造体
type CommentController struct {
	service services.CommentServicer
}

// コンストラクタ関数
func NewCommentController(s services.CommentServicer) *CommentController {
	return &CommentController{service: s}
}

// ハンドラメソッドの定義
// POST /comment
func (c *CommentController) PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	var reqComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&reqComment); err != nil {
		apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		http.Error(w, "fail to decode json\n", http.StatusBadRequest)
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		http.Error(w, "fail to post comment\n", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(comment)
}
