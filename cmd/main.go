package main

import (
	"database/sql"
	"ecom/cmd/api"
	"ecom/config"
	"ecom/db"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.NewMYSQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":3000", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	log.Println("Connecting to database")
	// log.Println("db details",)
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}
   
	log.Println("DB successfully connected!")
}
