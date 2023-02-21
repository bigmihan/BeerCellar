package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

func (bt *BeerTransport) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bt.writeLog(r.Method, r.URL, nil)
		w.Header().Add("Content-Type", "application/json")
		//fmt.Println(string(next))
		next.ServeHTTP(w, r)
	})

}

func (bt *BeerTransport) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag_getTokenFromRequest")
			return
		}

		email, err := bt.tokenMaker.GetEmailFromToken(token)
		if err != nil {
			bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag_GetEmailFromToken")
			return
		}
		if email == "" {
			bt.ResponseError(w, r.URL, errors.New("email is empty"), http.StatusBadRequest, "tag_email is empty")
			return

		}
		ctx := r.Context()

		ctx = context.WithValue(ctx, "email", email)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})

}

func getTokenFromRequest(r *http.Request) (string, error) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return "", errors.New("not found healer Authorization")
	}

	m := strings.Split(bearerToken, " ")
	if len(m) != 2 || m[0] != "Bearer" {
		return "", errors.New("header Authorization must be 'Bearer token'")
	}

	return m[1], nil
}

func (bt *BeerTransport) writeLog(description string, URL *url.URL, err error) {

	bt.logger.WriteLog(description, URL.Path, err)
}

func (bt *BeerTransport) ResponseError(w http.ResponseWriter, URL *url.URL, err error, status int, tagError string) {
	w.WriteHeader(status)

	m := map[string]string{"error": tagError + " " + err.Error()}
	jsByte, errMarshal := json.Marshal(m)
	if errMarshal != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(jsByte)
	}

	bt.writeLog("", URL, err)

}
