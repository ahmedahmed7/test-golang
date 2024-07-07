package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/entities"
	"github.com/japhy-tech/backend-test/repositories"
	"net/http"
	"strconv"
)

func LoadData(w http.ResponseWriter, r *http.Request) {
	err := repositories.LoadData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "data loaded"})
}
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
func AddPet(w http.ResponseWriter, r *http.Request) {
	var pet entities.Pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := repositories.AddPet(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pet)
}
func UpdatePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	var pet entities.Pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pet.ID = uint(id)

	if err := repositories.UpdatePet(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "pet updated"})
}

func FilterPets(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	weightMaleStr := r.URL.Query().Get("weightMale")
	weightFemaleStr := r.URL.Query().Get("weightFemale")

	var weightMale, weightFemale float64
	var err error
	if weightMaleStr != "" {
		weightMale, err = strconv.ParseFloat(weightMaleStr, 64)
		if err != nil {
			http.Error(w, "weightMale invalid", http.StatusBadRequest)
			return
		}
	}

	if weightFemaleStr != "" {
		weightFemale, err = strconv.ParseFloat(weightFemaleStr, 64)
		if err != nil {
			http.Error(w, "weightFemale invalid", http.StatusBadRequest)
			return
		}
	}
	pets, err := repositories.FilterPet(species, weightMale, weightFemale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pets)
}
