package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ShinnosukeSuzuki/practice-golang-api/api"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser     = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbDatabase = os.Getenv("DB_NAME")
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

func main() {
	// DB接続
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		log.Println("failed to connect database")
		return
	}

	// コントローラ型MyAppControllerのハンドラメソッドとパスを紐付け
	r := api.NewRouter(db)

	// シグナルを受け取るためのコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	log.Println("server start at :8080")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		// シグナルを受け取るまで待機
		<-ctx.Done()

		// 5秒のタイムアウト付きコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// サーバーをシャットダウン(新しい接続の受け付けを停止し、contextがキャンセルされたら終了する)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Shutdown(): %v", err)
		}
		defer wg.Done()
	}()

	// 正常にシャットダウンした場合はhttp.ErrServerClosedが返る
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	wg.Wait()

}
