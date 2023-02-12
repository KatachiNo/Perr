package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/KatachiNo/Perr/internal/user"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	handler := user.NewHandler()
	handler.Register(router)

	fmt.Print("Are u ready?\n")

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

	r := mux.NewRouter()

	log.Fatal(http.ListenAndServe(":8000", r))

}
