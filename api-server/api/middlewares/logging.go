package middlewares

import (
	"log"
	"net/http"
)

// 自作ResponseWriterを作成
type resLoggingWriter struct {
	http.ResponseWriter
	code int
}

// コンストラクトを作る
func NewResLoggingWriter(w http.ResponseWriter) *resLoggingWriter {
	return &resLoggingWriter{w, http.StatusOK}
}

// WriteHeaderをオーバーライド
func (w *resLoggingWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエスト情報をロギング
		log.Println(r.URL.Path, r.Method)

		// 自作のResponseWriterを作成
		rlw := NewResLoggingWriter(w)

		// ハンドラに渡す
		next.ServeHTTP(rlw, r)

		// 自作ResponseWriterからロギングしたいデータを出す
		log.Println("res: ", rlw.code)
	})
}
