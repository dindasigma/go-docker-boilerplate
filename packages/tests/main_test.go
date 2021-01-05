package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/datasources"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var userInstance = users.User{}

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
	datasources.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}
}

func refreshUserTable() error {
	err := datasources.DB.DropTableIfExists(&users.User{}).Error
	if err != nil {
		return err
	}
	err = datasources.DB.AutoMigrate(&users.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedUser() (users.User, error) {
	refreshUserTable()

	user := users.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "password",
	}

	err := datasources.DB.Model(&users.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() ([]users.User, error) {
	usersSeed := []users.User{
		users.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "password",
		},
		users.User{
			FirstName: "Doe",
			LastName:  "John",
			Email:     "doe@john.com",
			Password:  "password",
		},
	}

	for i, _ := range usersSeed {
		err := datasources.DB.Model(&users.User{}).Create(&usersSeed[i]).Error
		if err != nil {
			return []users.User{}, err
		}
	}
	return usersSeed, nil
}
