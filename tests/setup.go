package tests

import (
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/heinswanhtet/blogora-api/configs"
	"github.com/heinswanhtet/blogora-api/db"
	store "github.com/heinswanhtet/blogora-api/stores"
)

var s = SetupTesting()

func SetupTesting() *store.Store {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	database, err := db.Connect(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// db.Ping(database)

	s := store.NewStore(database)

	return s
}
