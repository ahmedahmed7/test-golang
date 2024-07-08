package main

import (
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
	//load data from csv
	r.HandleFunc("/loadData", controllers.LoadData).Methods(http.MethodGet)
	// CRUD
	r.HandleFunc("/getAllPets", controllers.GetAllHandler).Methods(http.MethodGet)
	r.HandleFunc("/getPet/{id}", controllers.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/deletePet/{id}", controllers.DeleteOne).Methods(http.MethodDelete)
	r.HandleFunc("/addPet", controllers.AddPet).Methods(http.MethodPost)
	r.HandleFunc("/updatePet/{id}", controllers.UpdatePet).Methods(http.MethodPatch)
	r.HandleFunc("/findByWeightSpecies", controllers.FilterPets).Methods(http.MethodGet)

	err = http.ListenAndServe(
		net.JoinHostPort("", ApiPort),
		r,
	)

	// =============================== Starting Msg ===============================
	utils.Logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))
}
