package config

import (
	"sync"

	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Type string `yaml:"type" env-default:"port"`
		Port string `yaml:"port" env-default:"8080"`
	} `yaml:"server"`

	Db struct {
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Dbname   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
		Password int    `yaml:"password" env:"PASSWORD" env-default:"1234"`
	} `yaml:"db"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		l := logg.GetLogger()
		l.Info("Getting configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("configs/config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)

			l.Info(help)
			l.Fatal(err)
		}
	})
	return instance
}
