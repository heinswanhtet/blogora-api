package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/heinswanhtet/blogora-api/utils"
)

func Connect(cfg *mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	ready := utils.GetColoredString("[ ready ]", utils.GREEN)
	log.Printf("%v Database connected successfully\n", ready)
}
