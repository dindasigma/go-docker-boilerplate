package api

import (
	"fmt"
	"log"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/application"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/datasources"
	"github.com/joho/godotenv"
)

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	datasources.InitializePostgres()

	application.Run(":8080")

}
