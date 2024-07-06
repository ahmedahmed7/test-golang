package controllers

import (
	"encoding/json"
	"github.com/japhy-tech/backend-test/repositories"
	"net/http"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	pets, err := repositories.GetAllPetsFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets)
}
