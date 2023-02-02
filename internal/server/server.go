package server

type Server struct {
	config *Config
}

func New(config *Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Start() error {
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
}
