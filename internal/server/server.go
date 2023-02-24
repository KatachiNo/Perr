package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

func newDB(dburl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dburl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// дальше код не нужен

//func (s *Server) configurationAPI() {
//r.HandleFunc("/login/authorization", authorization)
//r.HandleFunc("/login/registration", registration)
//r.HandleFunc("/login/changeData", changeData)
//r.HandleFunc("/login/forgetPassword", forgetPassword)
//
//r.HandleFunc("/products/all", getAllProducts).Methods("GET")
//r.Handle("/products/add", isAuthorized(addProduct)).Methods("POST")
//r.HandleFunc("/products/changeProductItem", changeProduct)
//r.HandleFunc("/products/delete", deleteProduct)
//
//r.HandleFunc("/products/getPriceStory", deleteProduct)
//r.HandleFunc("/products/cha", deleteProduct)
//s.router.HandleFunc("/test", s.hey())
//}

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

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	//var AllProducts = ConnectionWithDB("products")
	//json.NewEncoder(w).Encode(AllProducts)
}

func ConnectionWithDB(resp string, db *sql.DB) []ProductsTable {

	var Data []ProductsTable

	rows, _ := db.Query(`SELECT * FROM "products"`)
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

func changeData(w http.ResponseWriter, r *http.Request) {

}
func forgetPassword(w http.ResponseWriter, r *http.Request) {

}
func addProduct(w http.ResponseWriter, r *http.Request) {

}
func changeProduct(w http.ResponseWriter, r *http.Request) {

}
func deleteProduct(w http.ResponseWriter, r *http.Request) {

}
func registration(w http.ResponseWriter, r *http.Request) {

}
