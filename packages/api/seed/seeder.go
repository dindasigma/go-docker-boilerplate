package seed

import (
	"log"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"github.com/jinzhu/gorm"
)

var usersSeed = []users.User{
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

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&users.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&users.User{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i, _ := range usersSeed {
		err = db.Debug().Model(&users.User{}).Create(&usersSeed[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}
