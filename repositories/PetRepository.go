package repositories

import (
	"github.com/japhy-tech/backend-test/db"
	"github.com/japhy-tech/backend-test/entities"
)

func GetAllPetsFromDB() ([]entities.Pet, error) {
	rows, err := db.DB.Query("SELECT id, species,pet_size, name, average_male_adult_weight,average_female_adult_weight FROM pets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []entities.Pet
	for rows.Next() {
		var pet entities.Pet
		if err := rows.Scan(&pet.ID, &pet.Species, &pet.PetSize, &pet.Name, &pet.AverageMaleAdultWeight, &pet.AverageFemaleAdultWeight); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}
