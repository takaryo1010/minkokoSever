package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func changeNickname(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	newNickname := c.FormValue("newNickname")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// ユーザー名に一致するレコードを取得
	var storedPasswordHash []byte
	err = db.QueryRow("SELECT password_hash FROM Person WHERE username = ?", username).Scan(&storedPasswordHash)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword(storedPasswordHash, []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Incorrect password"})
	}

	// 新しいニックネームの長さを確認
	if len(newNickname) > 20 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nickname must be 20 characters or less"})
	}

	// ニックネームを更新
	_, err = db.Exec("UPDATE Person SET nickname = ? WHERE username = ?", newNickname, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating nickname"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Nickname updated successfully"})
}

func toggleAffiliationFlag(c echo.Context) error {
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

	// Affiliation テーブル内の exists_flag を切り替える
	_, err = db.Exec("UPDATE Affiliation SET exists_flag = NOT exists_flag WHERE username = ? AND location_id = ?", username, locationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error toggling affiliation flag"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Affiliation flag toggled successfully"})
}

// ユーザー名とパスワードを検証する関数
func authenticateUser(db *sql.DB, username, password string) error {
	var storedPasswordHash []byte
	err := db.QueryRow("SELECT password_hash FROM Person WHERE username = ?", username).Scan(&storedPasswordHash)
	if err == sql.ErrNoRows {
		return fmt.Errorf("User not found")
	} else if err != nil {
		return fmt.Errorf("Database query error")
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword(storedPasswordHash, []byte(password)); err != nil {
		return fmt.Errorf("Incorrect password")
	}

	return nil
}
func updateLocation(c echo.Context) error {
	locationID := c.FormValue("locationID")
	latitude := c.FormValue("latitude")
	longitude := c.FormValue("longitude")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Location テーブルのレコードを更新
	_, err = db.Exec("UPDATE Location SET latitude = ?, longitude = ? WHERE id = ?", latitude, longitude, locationID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating location"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Location updated successfully"})
}
