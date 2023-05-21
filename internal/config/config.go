package config

import (
	"path/filepath"
	"sync"

	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Type string `yaml:"type" env-default:"port"`
		Port string `yaml:"port" env-default:"8080"`
	} `yaml:"server"`

	PostgresDb `yaml:"postgresConfigDb"`
}

type PostgresDb struct {
	Username string `yaml:"username"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Dbname   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	// Password                 string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	MaxAttemptsForConnection string `yaml:"max-attempts-for-connection" env-default:"10"`
	MakeStartTables          string `yaml:"make-start-tables" env-default:"false"`
	MakeStartAdmin           string `yaml:"make-start-admin" env-default:"false"`
}

type EnvConf struct {
	Password        string `env:"PASSWORD" env:"PASSWORD" env-default:"postgres"`
	SigningKeyAdmin string `env:"SIGNINGKEYADMIN"`
	SigningKeyUser  string `env:"SIGNINGKEYUSER"`
}

var instance *Config
var instance2 *EnvConf
var once sync.Once

func GetConfig() (*Config, *EnvConf) {
	once.Do(func() {
		l := logg.GetLogger()
		l.Info("Getting configuration /configs/config.yml")
		instance = &Config{}

		path := filepath.Join("configs", "config.yml") //"configs/config.yml"
		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)

			l.Info(help)
			l.Fatal(err)
		}

		l.Info("Getting configuration .env")
		instance2 = &EnvConf{}

		if err := cleanenv.ReadConfig(".env", instance2); err != nil {
			help, _ := cleanenv.GetDescription(instance2, nil)

			l.Info(help)
			l.Fatal(err)
		}
	})
	return instance, instance2
}
