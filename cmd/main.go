package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KatachiNo/Perr/internal/postgresDataBase"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

type ProductsTable struct {
	Id              int
	Category        int
	QuantityOfGoods int
	Lastprice       decimal.Decimal
	PictureAddress  string
}

type CategoryTable struct {
	Id           int
	CategoryId   int
	CategoryName string
}

type ProductPriceStory struct {
	Id    int
	Date  string
	Price decimal.Decimal
}

type UserPassword struct {
	Username string
	Password string
}

func main() {
	fmt.Print("Test")

	r := mux.NewRouter()
	r.HandleFunc("/products/all", getAllProducts)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var AllProducts = ConnectionWithDB("Products")
	json.NewEncoder(w).Encode(AllProducts)
}

func ConnectionWithDB(resp string) []ProductsTable {
	//var Data []ProductsTable

	var connStr postgresDataBase.ConfigDB = postgresDataBase.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.host"),
	}

	db, _ := postgresDataBase.ConnectToDB(connStr)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	//rows, err := db.Query("SELECT * FROM Products")
	//for rows.Next() {
	//f

	//var t :=
	//	pkg.ConfigDB{
	//	Host     :viper.GetString("db.host"),
	//	Port     :viper.GetString("db.port"),
	//	Username :viper.GetString("db.username"),
	//	Password :os.Getenv("DB_PASSWORD"),
	//	DBName   :viper.GetString("db.dbname"),
	//	SSLMode  :viper.GetString("db.host"),
	//  }
	//
}
