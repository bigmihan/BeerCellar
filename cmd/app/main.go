package main

import (
	"BeerCellar/config"
	"BeerCellar/internal/Transports/rest"
	"BeerCellar/internal/logger"
	"BeerCellar/internal/repository/Pgsql"
	"BeerCellar/internal/service/BeerCellar"
	"BeerCellar/pkg/database"
	hasher2 "BeerCellar/pkg/hasher"
	"BeerCellar/pkg/tokenMaker/jwtToken"
	"log"

	"net/http"
	"time"
)

func main() {

	configData, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     configData.Host,
		Port:     configData.Port,
		Username: configData.Username,
		DBName:   configData.DBName,
		SSLMode:  configData.SSLMode,
		Password: configData.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := Pgsql.NewDataBase(db)

	Logger := logger.NewLogger(repository)
	userRepository := Pgsql.NewUsers(db)

	//var u = domain.UsersIn{
	//	Email:    "var@gmail.com",
	//	Password: "73616c744d696b65995f545174c4668d074a4b27f3aa39f42f0003938c8b5f92be5dd46a85f6e5ef",
	//}
	//b, err := userRepository.ExistUser(&u)
	//fmt.Println(b)
	//fmt.Println(err)
	//return

	hasher := hasher2.NewSha256Hash(configData.Salt)
	tokenMaker := jwtToken.NewToken(configData.Salt, configData.TokenExpirationMinute)

	service := BeerCellar.NewBeerService(repository)
	userService := BeerCellar.NewUserService(userRepository, hasher, tokenMaker)
	handler := rest.NewBeerTransport(service, userService, Logger, tokenMaker)

	// init & run server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
