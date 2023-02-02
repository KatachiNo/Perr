package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	pdb "github.com/KatachiNo/Perr/internal/postgresDataBase"
	"github.com/KatachiNo/Perr/internal/server"
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

var (
	serverConf string
)

func init() {
	flag.StringVar(&serverConf, "config-path", "configs/server.yaml", "path to config file")
}

func main() {
	flag.Parse()
	
	s := server.New(server.NewConfig())
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Are u ready?\n")

	var vspec = viper.New()
	vspec.AddConfigPath(".")
	vspec.SetConfigName("PerrAppE")
	vspec.SetConfigFile("env")
	//vspec.AutomaticEnv()
	e := vspec.ReadInConfig()

	if e != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", e))
	}

	t := vspec.GetInt("DB_PASSWORD")
	fmt.Printf("%s sd", t)

	r := mux.NewRouter()

	r.HandleFunc("/login/authorization", authorization)
	r.HandleFunc("/login/registration", registration)
	r.HandleFunc("/login/changeData", changeData)
	r.HandleFunc("/login/forgetPassword", forgetPassword)

	r.HandleFunc("/products/all", getAllProducts).Methods("GET")
	r.Handle("/products/add", isAuthorized(addProduct)).Methods("POST")
	r.HandleFunc("/products/changeProductItem", changeProduct)
	r.HandleFunc("/products/delete", deleteProduct)

	r.HandleFunc("/products/getPriceStory", deleteProduct)
	r.HandleFunc("/products/cha", deleteProduct)

	log.Fatal(http.ListenAndServe(":8000", r))

	// if err := serv.Run("8000"); err != nil {
	// 	log.Fatal(http.ListenAndServe(":8000", r))
	// }

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

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var userTest = User{
	Username: "1",
	Password: "1",
}

func authorization(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	json.NewDecoder(r.Body).Decode(&u)

}

func checkAuthorization(u User) string {
	if userTest.Username != u.Password && userTest.Password != u.Password {
		fmt.Print("is not ok")
	}

	validToken, err := generationOfTokenJWT()

	if err == nil {
		fmt.Print(err)
	}
	return validToken

}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

var mySigningKey = []byte("johenews")

func generationOfTokenJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Hour * 2160).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
	}

	return tokenString, nil
}

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var AllProducts = ConnectionWithDB("Products")
	json.NewEncoder(w).Encode(AllProducts)
}

func ConnectionWithDB(resp string) []ProductsTable {
	var Data []ProductsTable

	viper.SetConfigName("config")
	viper.AddConfigPath("./configs/")
	err := viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var connStr pdb.ConfigDB = pdb.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err333 := pdb.ConnectToDB(connStr)

	if err333 != nil {
		log.Fatalf(": %s", err333.Error())
	}

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
