package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Postgres struct {
	Host                  string
	Port                  int
	Username              string
	DBName                string
	SSLMode               string
	Password              string
	Salt                  string
	TokenExpirationMinute int
}

func New() (*Postgres, error) {
	postgr := new(Postgres)

	err := godotenv.Load()
	if err != nil {
		return postgr, err
	}

	postgr.Host = os.Getenv("Host")
	portInt, err := strconv.Atoi(os.Getenv("Port"))
	if err != nil {
		return postgr, err
	}
	TokenExpirationMinuteInt, err := strconv.Atoi(os.Getenv("TokenExpirationMinute"))
	if err != nil {
		return postgr, err
	}

	postgr.TokenExpirationMinute = TokenExpirationMinuteInt
	postgr.Port = portInt
	postgr.Username = os.Getenv("UsernameDB")
	postgr.DBName = os.Getenv("DBName")
	postgr.SSLMode = os.Getenv("SSLMode")
	postgr.Password = os.Getenv("Password")
	postgr.Salt = os.Getenv("Salt")

	return postgr, nil
}
