package datasources

import (
	"fmt"
	"log"
	"os"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/helpers"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/seed"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitializePostgres() {
	dbConfig := helpers.NewDatabase(
		os.Getenv("DB_MS"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_TIMEZONE"),
		os.Getenv("DB_SSL_MODE"),
		os.Getenv("DB_TIMEZONE"),
		"",
		"",
	)

	// Connect to DB
	var err error
	DB, err = dbConfig.Connect()

	if err != nil {
		log.Fatalf("Invalid db config: %v", err)
	} else {
		fmt.Printf("We are connected to the database")
	}

	if err != nil {
		panic("failed to connect database")
	}

	seed.Load(DB)
}
