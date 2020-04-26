package modeltests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/victorsteven/fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.FirstName, user.FirstName)
	assert.Equal(t, foundUser.LastName, user.LastName)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Password, user.Password)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID: 1,
		FirstName: "Johny",
		LastName: "Doe",
		Email: "johny@doe.com",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return // todo check return
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.FirstName, savedUser.FirstName)
	assert.Equal(t, newUser.LastName, savedUser.LastName)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Password, savedUser.Password)
}

func TestUpdateAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	updateUser := models.User{
		ID: 1,
		FirstName: "John",
		LastName: "Update",
		Email: "johny@update.com",
		Password: "password",
	}

	updatedUser, err := userUpdate.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	
	assert.Equal(t, updatedUser.ID, savedUser.ID)
	assert.Equal(t, updatedUser.FirstName, savedUser.FirstName)
	assert.Equal(t, updatedUser.LastName, savedUser.LastName)
	assert.Equal(t, updatedUser.Email, savedUser.Email)
	assert.Equal(t, updatedUser.Password, savedUser.Password)

}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}

	isDeleted, err := userInstance.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	assert.Equal(t, int(isDeleted), 1)
}