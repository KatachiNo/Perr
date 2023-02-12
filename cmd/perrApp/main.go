package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/KatachiNo/Perr/internal/user"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	handler := user.New()
	handler.Register(router)

	fmt.Print("Are u ready?\n")

	//uploadConfiguration()

	start(router)
}

func uploadConfiguration() {
	var vspec = viper.New()
	vspec.AddConfigPath(".")
	vspec.SetConfigName("PerrAppE")
	vspec.SetConfigFile("env")

	e := vspec.ReadInConfig()

	if e != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", e))
	}

	t := vspec.GetInt("DB_PASSWORD")
	fmt.Printf("%s sd", t)
}

func start(router *mux.Router) {

	serv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(serv.ListenAndServe())

}
