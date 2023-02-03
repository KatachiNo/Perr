package server

import (
	"io"
	"net/http"

	"github.com/KatachiNo/Perr/internal/postgresDataBase"
	"github.com/gorilla/mux"
)

type Server struct {
	config *Config
	router *mux.Router
	store  *postgresDataBase.Store
}

func New(config *Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	return nil

	s.configurationAPI()
	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.Openport, s.router)
}

func (s *Server) configureStore() error {
	st := postgresDataBase.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
}

func (s *Server) configurationAPI() {
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
	s.router.HandleFunc("/test", s.hey())
}

func (s *Server) hey() http.HandlerFunc {
	type request struct {
		name string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hey")
	}
}
