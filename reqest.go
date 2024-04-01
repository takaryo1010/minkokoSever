package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	Person struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
	}
)

func getPersonByID(c echo.Context) error {
	id := c.FormValue("id")
	fmt.Println(id)
	// データベース接続
	db, err := sql.Open("mysql", db_state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection error"})
	}
	defer db.Close()

	// PersonテーブルからIDに一致するレコードを取得
	var person Person
	err = db.QueryRow("SELECT id, name FROM Person WHERE id = ?", id).Scan(&person.Id, &person.Name)
	if err == sql.ErrNoRows {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Person not found"})
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
	}

	return c.JSON(http.StatusOK, person)
}
