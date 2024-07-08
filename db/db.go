package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/japhy-tech/backend-test/utils"
)

var DB *sql.DB

func InitDB(dsn string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		utils.Logger.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		utils.Logger.Fatal(err)
	}

	return DB, err
}
