package main

import (
	"net/http"
	"time"

	"github.com/KatachiNo/Perr/pkg/logg"
	_ "github.com/lib/pq"

	"github.com/KatachiNo/Perr/internal/user"
	"github.com/gorilla/mux"
)

func main() {
	l := logg.GetLogger()

	l.Info("create router")
	router := mux.NewRouter()

	l.Info("register user handler")
	handler := user.New()
	handler.Register(router)

	//uploadConfiguration()

	start(router)
}

//func uploadConfiguration() {
//	var vspec = viper.New()
//	vspec.AddConfigPath(".")
//	vspec.SetConfigName("PerrAppE")
//	vspec.SetConfigFile("env")
//
//	e := vspec.ReadInConfig()
//
//	if e != nil { // Handle errors reading the config file
//		panic(fmt.Errorf("fatal error config file: %w", e))
//	}
//
//	t := vspec.GetInt("DB_PASSWORD")
//	fmt.Printf("%s sd", t)
//}

func start(router *mux.Router) {
	l := logg.GetLogger()

	serv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	l.Info("server is listening port :8080")
	l.Fatal(serv.ListenAndServe())

}
