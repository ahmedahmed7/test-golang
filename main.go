package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/controllers"
	"github.com/japhy-tech/backend-test/database_actions"
	"github.com/japhy-tech/backend-test/db"
	"github.com/japhy-tech/backend-test/internal"
	"github.com/japhy-tech/backend-test/utils"
	"net"
	"net/http"
	"os"
	"strconv"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

func main() {

	err := database_actions.InitMigrator(MysqlDSN)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	db.InitDB(MysqlDSN)

	if err != nil {
		utils.Logger.Fatal(err.Error())
		os.Exit(1)
	}
	// Define the CSV file path
	csvFile := "./breeds.csv"

	// Open the CSV file
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
	// Insert data
	for i, record := range records {
		if i == 0 {
			// Skip the header row
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

	msg, err := database_actions.RunMigrate("up", 0)
	if err != nil {
		utils.Logger.Error(err.Error())
	} else {
		utils.Logger.Info(msg)
	}
	defer db.DB.Close()
	db.DB.SetMaxIdleConns(0)

	err = db.DB.Ping()
	if err != nil {
		utils.Logger.Fatal(err.Error())
		os.Exit(1)
	}

	utils.Logger.Info("Database connected")

	app := internal.NewApp(utils.Logger)

	r := mux.NewRouter()
	app.RegisterRoutes(r.PathPrefix("/v1").Subrouter())

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	r.HandleFunc("/getAllPets", controllers.GetAllHandler).Methods(http.MethodGet)
	r.HandleFunc("/getPet/{id}", controllers.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/deletePet/{id}", controllers.DeleteOne).Methods(http.MethodDelete)
	r.HandleFunc("/addPet", controllers.AddPet).Methods(http.MethodPost)
	r.HandleFunc("/updatePet/{id}", controllers.UpdatePet).Methods(http.MethodPatch)
	//r.HandleFunc("/findByWeightSpecies/{weight}/{species}", controllers.FindByWeightSpecies).Methods(http.MethodGet)

	err = http.ListenAndServe(
		net.JoinHostPort("", ApiPort),
		r,
	)

	// =============================== Starting Msg ===============================
	utils.Logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))
}
