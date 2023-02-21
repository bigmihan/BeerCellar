package jwtToken

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Token struct {
	salt             string
	expirationMinute int
}

func NewToken(salt string, expirationMinute int) *Token {

	return &Token{salt: salt, expirationMinute: expirationMinute}
}

func (t *Token) MakeToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"nbf":   time.Now().Add(time.Minute * time.Duration(t.expirationMinute)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(t.salt))

}

func (t *Token) GetEmailFromToken(tokenString string) (string, error) {
	//fmt.Println("filling Batches")
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjEyfQ.hft9uHiYiG5ZLIscBvWDvQ7eTvA4SC-7Tsuf3MQwrwQ"

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(t.salt), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims["email"]), nil
	} else {
		return "", errors.New("token not valid (tag token.Claims.(jwt.MapClaims); ok && token.Valid)")
	}
}
