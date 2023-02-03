package server

import "github.com/KatachiNo/Perr/internal/postgresDataBase"

type Config struct {
	Openport string `yaml:address`
	Store    *postgresDataBase.Config
}

// def
func NewConfig() *Config {
	return &Config{
		Openport: ":8080",
		Store: postgresDataBase.NewConfig(),
	}
}
