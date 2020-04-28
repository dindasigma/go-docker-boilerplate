package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/controllers"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models"
)

var server = controllers.Server{}
var userInstance = models.User{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}

	// connect to db before the test run
	databaseConnect()

	exitVal := m.Run()
	os.Exit(exitVal)
}

func databaseConnect() {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
	server.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedUser() (models.User, error) {
	refreshUserTable()

	user := models.User {
		FirstName: "John",
		LastName: "Doe",
		Email: "john@doe.com",
		Password: "password",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {
	users := []models.User{
		models.User{
			FirstName: "John",
			LastName: "Doe",
			Email: "john@doe.com",
			Password: "password",
		},
		models.User{
			FirstName: "Doe",
			LastName: "John",
			Email: "doe@john.com",
			Password: "password",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}