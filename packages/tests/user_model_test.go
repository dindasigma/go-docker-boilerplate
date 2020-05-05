package tests

import (
	"log"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models"
	"github.com/stretchr/testify/assert"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v\n", err)
		return
	}

	assert.Equal(t, 2, len(*users))
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

	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.FirstName, foundUser.FirstName)
	assert.Equal(t, user.LastName, foundUser.LastName)
	assert.Equal(t, user.Email, foundUser.Email)
	assert.Equal(t, user.Password, foundUser.Password)
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

	updatedUser, err := updateUser.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	
	assert.Equal(t, updateUser.ID, updatedUser.ID)
	assert.Equal(t, updateUser.FirstName, updatedUser.FirstName)
	assert.Equal(t, updateUser.LastName, updatedUser.LastName)
	assert.Equal(t, updateUser.Email, updatedUser.Email)
	assert.Equal(t, updateUser.Password, updatedUser.Password)

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

	assert.Equal(t, 1, int(isDeleted))
}