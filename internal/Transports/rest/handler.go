package rest

import (
	"BeerCellar/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type BeerService interface {
	InsertBeer(Id int, Name string, ABV float32, Quantity int) error
	TakeBeer(Name string, Quantity int, DaysAging int) ([]domain.CellarRemains, error)
	DeletePart(Name string, BatchString string) (int64, error)
	//WriteLog(description string, url string, err error)
}

type UserService interface {
	Create(user *domain.Users) error
	SignIn(user *domain.UsersIn) (string, error)
}

type LoggerService interface {
	WriteLog(description string, URL string, err error)
}

type BeerTransport struct {
	service    BeerService
	logger     LoggerService
	users      UserService
	tokenMaker domain.TokenMaker
}

func NewBeerTransport(service BeerService, user UserService, logger LoggerService, tokenMaker domain.TokenMaker) *BeerTransport {
	return &BeerTransport{service: service,
		logger:     logger,
		users:      user,
		tokenMaker: tokenMaker}
}

type inData struct {
	Id          int
	Name        string
	ABV         float32
	Quantity    int
	DaysAgain   int
	BatchString string
}

func (bt *BeerTransport) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(bt.loggerMiddleware)

	cellar := r.PathPrefix("/cellar").Subrouter()
	cellar.Use(bt.authMiddleware)
	cellar.HandleFunc("/insert", bt.HandlerInsert).Methods(http.MethodPost)
	cellar.HandleFunc("/take", bt.HandlerTake).Methods(http.MethodPost)
	cellar.HandleFunc("/delete", bt.HandlerDelete).Methods(http.MethodPost)

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/sign-up", bt.CreateUser).Methods(http.MethodPost)
	auth.HandleFunc("/sign-in", bt.SignIn).Methods(http.MethodGet)

	return r
}
