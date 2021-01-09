package seed

import (
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models/users"
	"gorm.io/gorm"
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
	// drop table if exist
	db.Debug().Migrator().DropTable(&users.User{})

	// Migrate the schema
	db.AutoMigrate(&users.User{})

	// Create
	for i, _ := range usersSeed {
		db.Create(&usersSeed[i])
	}
}
