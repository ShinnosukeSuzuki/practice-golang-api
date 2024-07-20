package common

import (
	"context"
	"net/http"
)

type TraceIDKeyType struct{}

func SetTraceID(ctx context.Context, traceID int) context.Context {
	// ctx に、(key: "traceID", value: 変数 traceID の値) をセット
	return context.WithValue(ctx, TraceIDKeyType{}, traceID)
}

func GetTraceID(ctx context.Context) int {
	id := ctx.Value(TraceIDKeyType{})
	if idInt, ok := id.(int); ok {
		return idInt
	}
	return 0
}

// コンテキストの中で name フィールドに対応させるキー構造体
type userNameKey struct{}

// コンテキストから name フィールドの値を取り出す関数
func GetUserName(ctx context.Context) string {
	id := ctx.Value(userNameKey{})

	if usernameStr, ok := id.(string); ok {
		return usernameStr
	}
	return ""
}

// コンテキストに name フィールドの値をセットする関数
func SetUserName(r *http.Request, name string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, userNameKey{}, name)
	return r.WithContext(ctx)
}
