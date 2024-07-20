package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/ShinnosukeSuzuki/practice-golang-api/apperrors"
	"github.com/ShinnosukeSuzuki/practice-golang-api/common"
	"google.golang.org/api/idtoken"
)

const (
	googleClientID = "[client_id]"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーから Authorization ヘッダーを取得
		authorization := r.Header.Get("Authorization")

		// AuthorizationフィールドがBearer [id_token]の形式になっているか検証
		authHeader := strings.Split(authorization, " ")
		if len(authHeader) != 2 {
			err := apperrors.RequiredAuthorizationHeader.Wrap(errors.New("invalid req header"), "invalid header")
			apperrors.ErrorHandler(w, r, err)
			return
		}

		bearer, idToken := authHeader[0], authHeader[1]
		if bearer != "Bearer" || idToken == "" {
			err := apperrors.RequiredAuthorizationHeader.Wrap(errors.New("invalid authorization header"), "invalid header")
			apperrors.ErrorHandler(w, r, err)
			return
		}

		// idTokenを検証する
		tokenValidator, err := idtoken.NewValidator(context.Background())
		if err != nil {
			err := apperrors.CannotMakeValidator.Wrap(err, "internal auth error")
			apperrors.ErrorHandler(w, r, err)
			return
		}

		payload, err := tokenValidator.Validate(context.Background(), idToken, googleClientID)
		if err != nil {
			err := apperrors.Unauthorizated.Wrap(err, "invalid id token")
			apperrors.ErrorHandler(w, r, err)
			return
		}

		// idTokenのペイロードからnameフィールドを取得
		name, ok := payload.Claims["name"]
		if !ok {
			err := apperrors.Unauthorizated.Wrap(err, "invalid id token")
			apperrors.ErrorHandler(w, r, err)
			return
		}

		// コンテキストにnameフィールドの値をセット
		r = common.SetUserName(r, name.(string))

		next.ServeHTTP(w, r)

	})
}
