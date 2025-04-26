package db

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import postgres driver
)

var DB *sqlx.DB

func Connect(dbUser, dbPassword, dbHost string, dbPort int, dbName string) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		strconv.Itoa(dbPort),
		dbName,
	)
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	log.Println("Connected to Postgres!")
}
