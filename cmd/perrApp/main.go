package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	pdb "github.com/KatachiNo/Perr/internal/postgresDataBase"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

type ProductsTable struct {
	Id              int             `json:"id"`
	Category        int             `json:"category"`
	QuantityOfGoods int             `json:"quantityofgoods"`
	Lastprice       decimal.Decimal `json:"lastprice"`
	PictureAddress  string          `json:"pictureaddress"`
}

type CategoryTable struct {
	Id           int    `json:"id"`
	CategoryId   int    `json:"categoryid"`
	CategoryName string `json:"categoryname"`
}

type ProductPriceStory struct {
	Id    int             `json:"id"`
	Date  string          `json:"date"`
	Price decimal.Decimal `json:"price"`
}

type UserPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	fmt.Print("Are u ready?\n")

	//serv := new(transport.Server)

	r := mux.NewRouter()
	r.HandleFunc("/products/all", getAllProducts)

	log.Fatal(http.ListenAndServe(":8000", r))

	// if err := serv.Run("8000"); err != nil {
	// 	log.Fatal(http.ListenAndServe(":8000", r))
	// }

}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var AllProducts = ConnectionWithDB("Products")
	json.NewEncoder(w).Encode(AllProducts)
}

func ConnectionWithDB(resp string) []ProductsTable {
	var Data []ProductsTable

	// viper.AddConfigPath("/configs")
	// viper.SetConfigName("config")
	// var connStr pdb.ConfigDB = pdb.ConfigDB{
	// 	Host:     viper.GetString("db.host"),
	// 	Port:     viper.GetString("db.port"),
	// 	Username: viper.GetString("db.username"),
	// 	Password: os.Getenv("DB_PASSWORD"),
	// 	DBName:   viper.GetString("db.dbname"),
	// 	SSLMode:  viper.GetString("db.host"),
	// }

	var connStr pdb.ConfigDB = pdb.ConfigDB{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "1234",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	db, err333 := pdb.ConnectToDB(connStr)

	if err333 != nil {
		log.Fatalf("error loading env variables: %s", err333.Error())
	}

	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("error loading env variables: %s", err.Error())
	//}

	rows, _ := db.Query(`SELECT * FROM "Products"`)
	for rows.Next() {
		var pTable = ProductsTable{}
		er := rows.Scan(&pTable.Category, &pTable.PictureAddress, &pTable.Id, &pTable.QuantityOfGoods, &pTable.Lastprice)

		if er != nil {
			log.Fatal(er)
		}
		Data = append(Data, pTable)
	}

	return Data
}
