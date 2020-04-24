package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models"
)

var users = []models.User{
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

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}