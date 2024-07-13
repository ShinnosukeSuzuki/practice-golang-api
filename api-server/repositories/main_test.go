package repositories_test

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// テスト全体で共有するsql.DBの変数
var testDB *sql.DB

var (
	dbUser     = "docker"
	dbPassword = "docker"
	dbDatabase = "sampledb"
	dbConn     = fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", dbUser, dbPassword, dbDatabase)
)

// DB接続する関数
func connectDB() error {
	var err error
	testDB, err = sql.Open("mysql", dbConn)
	if err != nil {
		return err
	}
	return nil
}

// テストデータのセットアップ
func setupTestData() error {
	// ./testdata/setupDB.sqlを実行してテストデータを初期化
	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-u", "docker", "sampledb", "--password=docker", "-e", "source ./testdata/setupDB.sql")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// DBのcleanup
func cleanupDB() error {
	// ./testdata/cleanupDB.sqlを実行してテストデータを初期化
	cmd := exec.Command("mysql", "-h", "127.0.0.1", "-u", "docker", "sampledb", "--password=docker", "-e", "source ./testdata/cleanupDB.sql")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 全テスト共通の前処理
func setup() error {
	if err := connectDB(); err != nil {
		return err
	}
	if err := cleanupDB(); err != nil {
		return err
	}
	if err := setupTestData(); err != nil {
		return err
	}
	return nil
}

// 全テスト共通の後処理
func teardown() {
	cleanupDB()
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
