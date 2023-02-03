package postgresDataBase

type Config struct{
	DatabaseURL string `yaml:database_url`
}

func NewConfig() *Config{
	return &Config{}
}
