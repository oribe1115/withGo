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

type CountryWithFewData struct {
	Code       string `json:"code,omitempty"`
	Name       string `json:"name,omitempty"`
	Population int    `json:"population,omitempty"`
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
	e.GET("/occupancy/:nameOfCity", Percentage)

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
	newCity := new(City)
	err := c.Bind(newCity)

	if err != nil {
		return c.JSON(http.StatusBadRequest, newCity)
	}

	db.Exec("INSERT INTO citesInJapan VALUES(?, ?, ?, ?, ?)", newCity.ID, newCity.Name, newCity.CountryCode, newCity.District, newCity.Population)

	return c.String(http.StatusOK, "Finished!")

}

func Percentage(c echo.Context) error {
	cityName := c.Param("nameOfCity")

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	thisCountry := CountryWithFewData{}

	db.Get(&thisCountry, "SELECT Code, Name, Population FROM country WHERE Code=?", city.CountryCode)

	return c.String(http.StatusOK, thisCountry.Name)

	occupaid := (city.Population / thisCountry.Population) * 100

	return c.String(http.StatusOK, string(occupaid)+"%")
}
