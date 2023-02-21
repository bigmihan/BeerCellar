package rest

import (
	"BeerCellar/internal/domain"
	"encoding/json"
	"io"
	"net/http"
)

func (bt *BeerTransport) CreateUser(w http.ResponseWriter, r *http.Request) {

	strByte, err := io.ReadAll(r.Body)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag4")
		return
	}

	user := new(domain.Users)
	err = json.Unmarshal(strByte, user)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag5")
		return
	}

	if err = user.Validate(); err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag6")
		return
	}
	err = bt.users.Create(user)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag7")
		return
	}

	w.WriteHeader(http.StatusOK)
	bt.writeLog("successfully:"+string(r.Method), r.URL, err)
}

func (bt *BeerTransport) SignIn(w http.ResponseWriter, r *http.Request) {
	strByte, err := io.ReadAll(r.Body)
	if err != nil {

		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag1")
		return
	}

	user := new(domain.UsersIn)
	err = json.Unmarshal(strByte, user)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag2")
		return
	}

	//if err = user.Validate(); err != nil {
	//	bt.ResponseError(w, r.URL, err, http.StatusBadRequest)
	//	return
	//}

	token, err := bt.users.SignIn(user)
	if err != nil {

		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag3")
		return
	}

	jsByte, errMarshal := json.Marshal(map[string]string{"token": token})
	if errMarshal != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(jsByte)
	}

	w.WriteHeader(http.StatusOK)
	bt.writeLog("successfully:"+string(r.Method), r.URL, err)

}
