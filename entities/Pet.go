package entities

type Pet struct {
	ID                       uint    `json:"id"`
	Species                  string  `json:"species"`
	PetSize                  string  `json:"pet_size"`
	Name                     string  `json:"name"`
	AverageMaleAdultWeight   float64 `json:"average_male_adult_weight"`
	AverageFemaleAdultWeight float64 `json:"average_female_adult_weight"`
}
