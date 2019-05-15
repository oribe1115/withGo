package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db
	e := echo.New()

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.POST("/addingCity", AddNewCity)

	e.Start(":10200")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	if city.Name == "" {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func AddNewCity(c echo.Context) error {

}
