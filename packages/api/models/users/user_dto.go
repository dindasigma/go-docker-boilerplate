package users

import "time"

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id" example:"1"`
	FirstName string    `gorm:"size:255;not null;" json:"first_name" example:"John"`
	LastName  string    `gorm:"size:255;" json:"last_name" example:"Doe"`
	Email     string    `gorm:"size:255;not null;unique" json:"email" example:"john@doe.com"`
	Password  string    `gorm:"size:100;not null;" json:"password" example:"password"`
	Role      string    `gorm:"size:100;" json:"role" example:"admin"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" example:"2020-09-06T15:17:17.769031568Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2020-09-06T15:17:17.769031568Z"`
}
