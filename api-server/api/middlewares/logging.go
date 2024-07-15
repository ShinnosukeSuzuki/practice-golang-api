package middlewares

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエスト情報をロギング
		log.Println(r.URL.Path, r.Method)
		next.ServeHTTP(w, r)
	})
}
