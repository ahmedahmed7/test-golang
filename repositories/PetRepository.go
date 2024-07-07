package repositories

import (
	"errors"
	"github.com/japhy-tech/backend-test/db"
	"github.com/japhy-tech/backend-test/entities"
)

func GetAllPets() ([]entities.Pet, error) {
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
func GetPetById(id int) (*entities.Pet, error) {
	query := "SELECT id, species,pet_size, name, average_male_adult_weight,average_female_adult_weight FROM pets where id=?"
	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pet entities.Pet

	if rows.Next() {
		if err := rows.Scan(&pet.ID, &pet.Species, &pet.PetSize, &pet.Name, &pet.AverageMaleAdultWeight, &pet.AverageFemaleAdultWeight); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("pet not found")
	}

	return &pet, nil
}
func DeletePetById(id int) error {
	query := "DELETE FROM pets WHERE id=?"
	rows, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("pet doesn't exist")
	}

	return nil
}
func AddPet(pet *entities.Pet) error {
	query := "INSERT INTO pets (species, pet_size, name, average_male_adult_weight, average_female_adult_weight) VALUES (?, ?, ?, ?, ?)"
	rows, err := db.DB.Exec(query, pet.Species, pet.PetSize, pet.Name, pet.AverageMaleAdultWeight, pet.AverageFemaleAdultWeight)
	if err != nil {
		return err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return err
	}

	pet.ID = uint(id)

	return nil
}
func UpdatePet(pet *entities.Pet) error {
	query := `
        UPDATE pets 
        SET species = ?, pet_size = ?, name = ?, average_male_adult_weight = ?, average_female_adult_weight = ?
        WHERE id = ?
    `
	_, err := db.DB.Exec(query, pet.Species, pet.PetSize, pet.Name, pet.AverageMaleAdultWeight, pet.AverageFemaleAdultWeight, pet.ID)
	if err != nil {
		return err
	}
	return nil
}
