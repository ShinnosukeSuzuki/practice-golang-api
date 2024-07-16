package middlewares

import (
	"context"
	"sync"
)

var (
	logNo int = 1
	mu    sync.Mutex
)

func newTraceID() int {
	var no int

	mu.Lock()
	no = logNo
	logNo += 1
	mu.Unlock()

	return no
}

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
