package cmd

import (
	"encoding/json"
	"net/http"
)

var couter int

func (api objectAPIHandlers) GetImageHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var obj struct {
		Ss string
		Id int
	}
	obj.Ss = "dsfdsf"
	couter++
	obj.Id = couter
	json.NewEncoder(w).Encode(obj)
}
