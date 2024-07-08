package repositories

import (
	"encoding/csv"
	"errors"
	"github.com/japhy-tech/backend-test/db"
	"github.com/japhy-tech/backend-test/entities"
	"github.com/japhy-tech/backend-test/utils"
	"os"
	"strconv"
)

func LoadData() error {
	csvFile := "./breeds.csv"

	file, err := os.Open(csvFile)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	stmt, err := db.DB.Prepare(`
				        INSERT INTO pets (id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight)
				        VALUES (?, ?, ?, ?, ?, ?)
				         ON DUPLICATE KEY UPDATE
		            species = VALUES(species),
		            pet_size = VALUES(pet_size),
		            name = VALUES(name),
		            average_male_adult_weight = VALUES(average_male_adult_weight),
		            average_female_adult_weight = VALUES(average_female_adult_weight)

				    `)
	defer stmt.Close()
	for i, record := range records {
		if i == 0 {
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
		avgMaleWeight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
		avgFemaleWeight, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
		_, err = stmt.Exec(id, record[1], record[2], record[3], avgMaleWeight, avgFemaleWeight)
		if err != nil {
			utils.Logger.Fatal(err.Error())
		}
	}
	return err
}
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

func FilterPet(species string, weightMale float64, weightFemale float64) ([]entities.Pet, error) {
	query := "SELECT id, species,pet_size, name, average_male_adult_weight,average_female_adult_weight FROM pets where species=? and (average_male_adult_weight=? or average_female_adult_weight=?) "
	rows, err := db.DB.Query(query, species, weightMale, weightFemale)
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
