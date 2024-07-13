package repositories_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// テスト全体で共有するsql.DBの変数
var testDB *sql.DB

// 全テスト共通の前処理
func setup() error {
	// DBへの接続
	dbUser := "docker"
	dbPassword := "docker"
	dbDatabase := "sampledb"
	dbConn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
	var err error
	testDB, err = sql.Open("mysql", dbConn)
	if err != nil {
		// DB接続エラーの場合はテストを中止する
		return err
	}
	return nil
}

// 全テスト共通の後処理
func teardown() {
	testDB.Close()
}

func TestMain(m *testing.M) {
	// テスト全体の前処理
	err := setup()
	if err != nil {
		os.Exit(1)
	}

	// テストの実行
	m.Run()

	// テスト全体の後処理
	teardown()
}
