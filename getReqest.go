package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getPersonByLocationID(c echo.Context) error {
	id := c.FormValue("id")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Affiliation テーブルと Person テーブルを結合して、指定された Location の ID に関連するニックネームと存在情報を取得
	rows, err := db.Query("SELECT p.nickname, a.exists_flag FROM Person p JOIN Affiliation a ON p.username = a.username WHERE a.location_id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}
	defer rows.Close()

	// ニックネームと存在情報のペアを格納するマップ
	personInfo := make(map[string]bool)

	// 取得したデータをマップに追加
	for rows.Next() {
		var nickname string
		var existsFlag bool
		err := rows.Scan(&nickname, &existsFlag)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning rows"})
		}
		personInfo[nickname] = existsFlag
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error iterating over rows"})
	}

	return c.JSON(http.StatusOK, personInfo)
}

func getLocationInfo(c echo.Context) error {
	locationID := c.QueryParam("locationID")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Location テーブルから LocationID に一致するレコードを取得
	var location Location
	err = db.QueryRow("SELECT id, name, latitude, longitude FROM Location WHERE id = ?", locationID).Scan(&location.ID, &location.Name, &location.Latitude, &location.Longitude)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Location not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}

	// Location 構造体を JSON 形式に変換して返す
	response, err := json.Marshal(location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error marshaling JSON"})
	}

	return c.JSON(http.StatusOK, response)
}
func getLocationsByUsername(c echo.Context) error {
	username := c.QueryParam("username")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Person テーブルから指定されたユーザーネームに関連する Location を取得するクエリ
	query := `SELECT l.id, l.name, l.latitude, l.longitude
			  FROM Location l
			  JOIN Affiliation a ON l.id = a.location_id
			  WHERE a.username = ?`

	rows, err := db.Query(query, username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}
	defer rows.Close()

	// 取得した Location を格納するスライス
	locations := make([]Location, 0)

	// 取得した Location をスライスに追加
	for rows.Next() {
		var location Location
		err := rows.Scan(&location.ID, &location.Name, &location.Latitude, &location.Longitude)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning rows"})
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error iterating over rows"})
	}
	return c.JSON(http.StatusOK, locations)
}
