package server

type Config struct {
	openport string `yaml:address`
}


//def
func NewConfig() *Config {
	return &Config{
		openport: ":8080",
	}
}
