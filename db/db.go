package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var Db *sql.DB
var err error

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	Db, err = sql.Open("mysql", os.Getenv("DNS"))
	if err != nil {
		panic(err.Error())
	}
}
