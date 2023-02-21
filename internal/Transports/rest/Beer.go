package rest

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"net/http"
)

func (bt *BeerTransport) HandlerTake(w http.ResponseWriter, r *http.Request) {

	strByte, err := io.ReadAll(r.Body)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag8")
		return
	}

	part := new(inData)
	err = json.Unmarshal(strByte, part)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag9")
		return
	}

	s, err := bt.service.TakeBeer(part.Name, part.Quantity, part.DaysAgain)

	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag10")
		return
	}

	sJson, err := json.Marshal(s)
	if err != nil {

		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag11")
		return
	}

	w.Write(sJson)
	w.WriteHeader(http.StatusOK)
	bt.writeLog("successfully:"+string(r.Method), r.URL, err)
}
func (bt *BeerTransport) HandlerInsert(w http.ResponseWriter, r *http.Request) {

	part := new(inData)
	strByte, err := io.ReadAll(r.Body)
	if err != nil {

		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag12")
		return
	}
	json.Unmarshal(strByte, part)
	fmt.Println(part)
	err = bt.service.InsertBeer(part.Id, part.Name, part.ABV, part.Quantity)

	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag13")
		return
	}
	bt.writeLog("successfully:"+string(r.Method), r.URL, err)
	w.WriteHeader(http.StatusOK)
}

func (bt *BeerTransport) HandlerDelete(w http.ResponseWriter, r *http.Request) {

	part := new(inData)
	strByte, err := io.ReadAll(r.Body)
	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusBadRequest, "tag14")
		return
	}
	json.Unmarshal(strByte, part)

	s, err := bt.service.DeletePart(part.Name, part.BatchString)

	if err != nil {
		bt.ResponseError(w, r.URL, err, http.StatusInternalServerError, "tag15")
		return
	}

	w.Write([]byte(fmt.Sprintf("count string delete %d", s)))
	bt.writeLog(fmt.Sprintf("successfully:%s count string delete %d", string(r.Method), s), r.URL, err)
	w.WriteHeader(http.StatusOK)
}
