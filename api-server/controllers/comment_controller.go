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
		err = apperrors.ReqBodyDecodeFailed.Wrap(err, "bad request body")
		apperrors.ErrorHandler(w, r, err)
	}

	comment, err := c.service.PostCommentService(reqComment)
	if err != nil {
		apperrors.ErrorHandler(w, r, err)
		return
	}
	json.NewEncoder(w).Encode(comment)
}
