package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func addPerson(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// パスワードをハッシュ化
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(passwordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
	}

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// ユーザー名が既に存在するかチェック
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Person WHERE username = ?", username).Scan(&count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}
	if count > 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username already exists"})
	}

	// Person テーブルにレコードを挿入
	_, err = db.Exec("INSERT INTO Person (username, password_hash) VALUES (?, ?)", username, passwordHash)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error inserting into database"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Person added successfully"})
}
func addAffiliation(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	locationID := c.FormValue("locationID")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// ユーザー名とパスワードの検証
	if err := authenticateUser(db, username, password); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Affiliation テーブルにレコードを挿入
	_, err = db.Exec("INSERT INTO Affiliation (username, location_id, exists_flag) VALUES (?, ?, false)", username, locationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error adding affiliation"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Affiliation added successfully"})
}

func addLocation(c echo.Context) error {
	name := c.FormValue("name")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// 新しい Location ID を生成して重複チェック
	locationID, err := generateUniqueLocationID(db)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error generating location ID"})
	}

	// Location テーブルにレコードを挿入
	_, err = db.Exec("INSERT INTO Location (id, name, latitude, longitude) VALUES (?, ?, 0.0, 0.0)", locationID, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error adding location"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Location added successfully"})
}

// 重複のない Location ID を生成する関数
func generateUniqueLocationID(db *sql.DB) (string, error) {
	rand.Seed(time.Now().UnixNano())

	for {
		// ランダムな5桁以上の数字と英小文字の組み合わせを生成
		id := ""
		for i := 0; i < 5; i++ {
			if rand.Intn(2) == 0 {
				id += fmt.Sprintf("%c", rand.Intn(10)+48)
			} else {
				id += fmt.Sprintf("%c", rand.Intn(27)+97) // 英小文字
			}
		}

		// 重複チェック
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM Location WHERE id = ?", id).Scan(&count)
		if err != nil {
			return "", err
		}
		if count == 0 {
			return id, nil
		}
	}
}
