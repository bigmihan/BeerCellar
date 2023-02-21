package BeerCellar

import (
	"BeerCellar/internal/domain"
	"errors"
)

type userDataBase interface {
	Create(user *domain.Users) error
	ExistUser(user *domain.UsersIn) (bool, error)
}

type hasher interface {
	Hash(password string) string
}

type UserService struct {
	UserRepository userDataBase
	Hasher         hasher
	tokenMaker     domain.TokenMaker
}

func NewUserService(UserRepository userDataBase, Hasher hasher, TokenMaker domain.TokenMaker) *UserService {
	return &UserService{
		UserRepository: UserRepository,
		Hasher:         Hasher,
		tokenMaker:     TokenMaker}
}

func (userService *UserService) Create(user *domain.Users) error {

	user.Password = userService.Hasher.Hash(user.Password)
	return userService.UserRepository.Create(user)
}

func (userService *UserService) SignIn(user *domain.UsersIn) (string, error) {

	user.Password = userService.Hasher.Hash(user.Password)
	ExistUser, err := userService.UserRepository.ExistUser(user)
	if err != nil {
		return "", err
	}
	if !ExistUser {
		return "", errors.New("user not exist")
	}

	return userService.tokenMaker.MakeToken(user.Email)
}
