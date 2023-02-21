package Pgsql

import (
	"BeerCellar/internal/domain"
	"database/sql"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (Udb *Users) Create(user *domain.Users) error {
	_, err := Udb.db.Exec("insert into users (email, password, name) values ($1,$2,$3)", user.Email, user.Password, user.Name)
	return err
}

func (Udb *Users) ExistUser(user *domain.UsersIn) (bool, error) {
	rows, err := Udb.db.Query("select * from users where email = $1 and password = $2", user.Email, user.Password)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}
