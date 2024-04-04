package main

import "database/sql"

type Location struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Person struct {
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	PasswordHash []byte `json:"-"`
}
type NicknameAndExist struct {
	Username  string `json:"-"`
	Nickname   string `json:"nickname"`
	Exist_flag bool           `json:"exists_flag"`
}
type NicknameAndExistAndUsername struct {
	Username  string `json:"username"`
	Nickname   sql.NullString `json:"nickname"`
	Exist_flag bool           `json:"exists_flag"`
}
