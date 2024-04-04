package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()
	fmt.Println("------------------------")
	// ユーザー名に一致するパスワードハッシュを取得
	var storedPasswordHash []byte
	err = db.QueryRow("SELECT password_hash FROM Person WHERE username = ?", username).Scan(&storedPasswordHash)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}
	fmt.Println(password)
	fmt.Println("------------------------")
	// パスワードの検証

	// パスワードの検証
	err = bcrypt.CompareHashAndPassword(storedPasswordHash, []byte(password))
	if err != nil {

		fmt.Println(storedPasswordHash)
		fmt.Println("------------------------")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Login successful"})
}

func getPersonByLocationID(c echo.Context) error {
	id := c.FormValue("id")
fmt.Println("aaaaaaaaaaaaaaaaaaaaaa"+id)
	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Affiliation テーブルと Person テーブルを結合して、指定された Location の ID に関連するニックネームと存在情報を取得
	rows, err := db.Query("SELECT a.username, p.nickname, a.exists_flag FROM Person p JOIN Affiliation a ON p.username = a.username WHERE a.location_id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}
	defer rows.Close()

	// ニックネームと存在情報のペアを格納するマップ
	var personInfo []NicknameAndExist
	personInfo = nil 

	// 取得したデータをマップに追加
	for rows.Next() {
		var info NicknameAndExistAndUsername
		var res NicknameAndExist
		err := rows.Scan(&info.Username,&info.Nickname, &info.Exist_flag)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error scanning rows"})
		}
		if info.Nickname.Valid{
			res.Nickname = info.Nickname.String
			res.Exist_flag = info.Exist_flag
			personInfo = append(personInfo,res)
		}else{
			res.Nickname = info.Username
			res.Exist_flag = info.Exist_flag
			personInfo = append(personInfo,res)
			}
		
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
	username := c.FormValue("username")

	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// Person テーブルから指定されたユーザーネームに関連する Location を取得するクエリ
	query := "SELECT Location.* FROM Affiliation JOIN Location ON Affiliation.location_id = Location.id WHERE Affiliation.username = ?"
	// SQLの準備
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return err // エラーを返す
	}
	defer stmt.Close() // ステートメントを閉じる
	rows, err :=stmt.Query(username)
	fmt.Println(username)
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
		fmt.Println(location)
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error iterating over rows"})
	}
	return c.JSON(http.StatusOK, locations)
}
