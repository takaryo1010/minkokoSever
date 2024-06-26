package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func connectOnly() string {
	// データベースのハンドルを取得する
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		// ここではエラーを返さない
		log.Fatal(err)
	}
	defer db.Close()

	// 実際に接続する
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	} else {
		log.Println("データベース接続完了")
		return "データベース接続完了"
	}
}

// func sqlInsert(name string) {
// 	// データベースのハンドルを取得する
// 	db, err := sql.Open("mysql", db_state)
// 	if err != nil {
// 		// ここではエラーを返さない
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// SQLの準備
// 	ins, err := db.Prepare("INSERT INTO Person VALUES(?,?)")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ins.Close()

// 	// SQLの実行
// 	res, err := ins.Exec(0, name)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// 結果の取得
// 	lastInsertID, err := res.LastInsertId()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println(lastInsertID)
// }
