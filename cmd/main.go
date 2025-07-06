package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/heinswanhtet/blogora-api/cmd/api"
	"github.com/heinswanhtet/blogora-api/configs"
	"github.com/heinswanhtet/blogora-api/db"
)

func main() {
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.UTC,
	}

	database, err := db.Connect(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	db.Ping(database)

	server := api.NewAPIServer(
		fmt.Sprintf("%s:%s", configs.Envs.Host, configs.Envs.Port),
		database,
	)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
