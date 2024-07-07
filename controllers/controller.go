package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/repositories"
	"net/http"
	"strconv"
)

func GetAllHandler(w http.ResponseWriter, r *http.Request) {
	pets, err := repositories.GetAllPets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pets)
}
func GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	pet, err := repositories.GetPetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pet)
}
func DeleteOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}
	err = repositories.DeletePetById(id)
	if err != nil {
		if err.Error() == "pet doesn't exist" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "pet deleted"})
}
