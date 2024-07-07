package main

import (
	"encoding/csv"
	"fmt"
	"github.com/japhy-tech/backend-test/controllers"
	"github.com/japhy-tech/backend-test/db"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/database_actions"
	"github.com/japhy-tech/backend-test/internal"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

func main() {
	logger := charmLog.NewWithOptions(os.Stderr, charmLog.Options{
		Formatter:       charmLog.TextFormatter,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "🧑‍💻 backend-test",
		Level:           charmLog.DebugLevel,
	})

	err := database_actions.InitMigrator(MysqlDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}
	db.InitDB(MysqlDSN)

	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	// Define the CSV file path
	csvFile := "./breeds.csv"

	// Open the CSV file
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
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
	// Insert data into the MySQL table
	for i, record := range records {
		if i == 0 {
			// Skip the header row
			continue
		}
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("Failed to convert id to integer: %v", err)
		}
		avgMaleWeight, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Fatalf("Failed to convert average_male_adult_weight to float: %v", err)
		}
		avgFemaleWeight, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Fatalf("Failed to convert average_female_adult_weight to float: %v", err)
		}
		_, err = stmt.Exec(id, record[1], record[2], record[3], avgMaleWeight, avgFemaleWeight)
		if err != nil {
			log.Fatalf("Failed to execute statement: %v", err)
		}
	}

	fmt.Println("Data inserted successfully")
	msg, err := database_actions.RunMigrate("up", 0)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info(msg)
	}
	defer db.DB.Close()
	db.DB.SetMaxIdleConns(0)

	err = db.DB.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	logger.Info("Database connected")

	app := internal.NewApp(logger)

	r := mux.NewRouter()
	app.RegisterRoutes(r.PathPrefix("/v1").Subrouter())

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	r.HandleFunc("/getAllPets", controllers.GetAllHandler).Methods(http.MethodGet)
	r.HandleFunc("/getPet/{id}", controllers.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/deletePet/{id}", controllers.DeleteOne).Methods(http.MethodDelete)

	err = http.ListenAndServe(
		net.JoinHostPort("", ApiPort),
		r,
	)

	// =============================== Starting Msg ===============================
	logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))
}
