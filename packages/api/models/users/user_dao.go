package users

import (
	"context"
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/crypto"
	"gorm.io/gorm"
)

func (u *User) beforeSave() error {
	hashedPassword, err := crypto.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Role = html.EscapeString(strings.TrimSpace(u.Role))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.FirstName == "" {
			return errors.New("Required First Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.FirstName == "" {
			return errors.New("Required First Name")
		}
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) Save(ctx context.Context, db *gorm.DB) (*User, error) {
	err := db.WithContext(ctx).Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAll(ctx context.Context, db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.WithContext(ctx).Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindByID(ctx context.Context, db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.WithContext(ctx).Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}

	if err == gorm.ErrRecordNotFound {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) Update(ctx context.Context, db *gorm.DB, uid uint32) (*User, error) {
	// To hash the password
	err := u.beforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.WithContext(ctx).Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"role":       u.Role,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.WithContext(ctx).Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) Delete(ctx context.Context, db *gorm.DB, uid uint32) (int64, error) {
	db = db.WithContext(ctx).Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) Check(db *gorm.DB, email string) error {

	err := db.Debug().Model(User{}).Where("email = ?", email).Take(&u).Error
	if err != nil {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		return errors.New("User Not Found")
	}
	return nil
}
