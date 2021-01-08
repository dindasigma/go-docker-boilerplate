package datasources

import (
	"fmt"
	"log"
	"os"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/helpers"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/seed"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
		"",
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

	DB.Debug().AutoMigrate(&users.User{})
	seed.Load(DB)
}
