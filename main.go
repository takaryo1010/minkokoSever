package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db_state string

func main() {
	// ファイルからdb_stateを読み込む
	dbStateFromFile, err := os.ReadFile("db_state.txt")
	if err != nil {
		panic(err)
	}
	db_state = string(dbStateFromFile)

	// インスタンスを作成
	e := echo.New()

	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルートを設定
	e.GET("/aa", hello) // ローカル環境の場合、http://localhost:1323/ にGETアクセスされるとhelloハンドラーを実行する
	e.POST("/getPersonByLocationID", getPersonByLocationID)//id
	e.GET("/", connect_check)
	e.POST("/login",login)//username password
	e.POST("/addPerson",addPerson)// username password
	e.POST("/addAffiliation",addAffiliation)// username password locationID
	e.POST("/addLocation",addLocation)//name
	e.POST("/getLocationInfo",getLocationInfo);//locationID
	e.POST("/getLocationsByUsername",getLocationsByUsername);//username
	e.POST("/changeNickname",changeNickname);//username password newNickname
	e.POST("/toggleAffiliationFlag",toggleAffiliationFlag)//username password locationID
	e.POST("/updateLocation",updateLocation)//locationID latitude longitude




	e.Logger.Fatal(e.Start(":1323"))
	fmt.Println(os.Getenv("PORT"))
}

// ハンドラーを定義
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func connect_check(c echo.Context) error {
	res := connectOnly()
	return c.String(http.StatusOK, res)
}
