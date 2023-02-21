package domain

import "github.com/go-playground/validator/v10"

type Users struct {
	Email    string `validate:"email"`
	Password string `validate:"gte=8"`
	Name     string `validate:"required"`
}

type UsersIn struct {
	Email    string `validate:"email"`
	Password string `validate:"gte=3"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (u *Users) Validate() error {
	return validate.Struct(u)
}

func (u *UsersIn) Validate() error {
	return validate.Struct(u)
}

type TokenMaker interface {
	MakeToken(email string) (string, error)
	GetEmailFromToken(tokenString string) (string, error)
}
