package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Ok bool `json:"ok"`
}

func HealtzHandler(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(Response{ Ok: true })
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
