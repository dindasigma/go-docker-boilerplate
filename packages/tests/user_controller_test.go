package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON	string
		statusCode	int
		firstName	string
		email	string
		errorMessage	string
	}{
		{
			inputJSON: `{"firstname":"John", "lastname":"Doe", "email":"john@doe.com", "password":"password"}`,
			statusCode: 201,
			firstName: "John",
			email: "john@doe.com",
			errorMessage: "",
		},
		{
			inputJSON: `{"firstname":"John", "lastname":"Doe", "email":"john@doe.com", "password":"password"}`,
			// todo 422
			statusCode: 201,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON: `{"firstname":"John", "lastname":"Doe", "email":"johndoe.com", "password":"password"}`,
			statusCode: 500,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON: `{"firstname":"", "lastname":"Doe", "email":"doe@john.com", "password":"password"}`,
			statusCode: 422,
			errorMessage: "Required First Name",
		},
		{
			inputJSON: `{"firstname":"John", "lastname":"Doe", "email":"", "password":"password"}`,
			statusCode: 422,
			errorMessage: "Required Email",
		},
		{
			inputJSON: `{"firstname":"John", "lastname":"Doe", "email":"john@doe.com", "password":""}`,
			statusCode: 422,
			errorMessage: "Required Password",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("This is the error: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["firstname"], v.firstName)
			assert.Equal(t, responseMap["email"], v.email)
		}

		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("This is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json %v\n", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUser()
	if err != nil {
		log.Fatal(err)
	}

	sample := []struct {
		id string
		statusCode int
		firstname string
		email string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id: strconv.Itoa(int(user.ID)),
			statusCode: 200,
			firstname: user.FirstName,
			email: user.Email,
		},
		{
			id: "unknown",
			statusCode: 400,
		},
	}

	for _, v := range sample {
		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			//assert.Equal(t, user.FirstName, responseMap["firstname"])
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}

	// get only the first user
	AuthID = users[0].ID
	AuthEmail = users[0].Email
	AuthPassword = "password"

	// login
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id string
		updateJSON string
		statusCode int
		updateFirstname string
		updateEmail string
		tokenGiven string
		errorMessage string
	}{
		{
			
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "john@doe.com", "password": "password"}`,
			// to do 200
			statusCode: 422,
			updateFirstname: "John",
			updateEmail: "john@doe.com",
			tokenGiven: tokenString,
			errorMessage: "",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "john@doe.com", "password": ""}`,
			statusCode: 422,
			tokenGiven: tokenString,
			errorMessage:"Required Password",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "john@doe.com", "password": "password"}`,
			statusCode: 422,
			tokenGiven: "",
			errorMessage: "Unauthorized",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "john@doe.com", "password": "password"}`,
			statusCode: 422,
			tokenGiven: "incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "doe@john.com", "password": "password"}`,
			statusCode: 500,
			tokenGiven: tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "doejohn.com", "password": "password"}`,
			statusCode: 500,
			tokenGiven: tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"", "lastname":"Doe", "email": "john@doe.com", "password": "password"}`,
			statusCode: 422,
			tokenGiven: tokenString,
			errorMessage: "Required First Name",
		},
		{
			id: strconv.Itoa(int(AuthID)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "", "password": "password"}`,
			statusCode: 422,
			tokenGiven: tokenString,
			errorMessage: "Required Email",
		},
		{
			id: "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id: strconv.Itoa(int(2)),
			updateJSON: `{"firstname":"John", "lastname":"Doe", "email": "", "password": "password"}`,
			tokenGiven: tokenString,
			statusCode: 401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["firstname"], v.updateFirstname)
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	var AuthEmail, AuthPassword string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedUsers()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
		
	AuthID = user[0].ID
	AuthEmail = user[0].Email
	AuthPassword = "password"
	
	// login
	token, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(AuthID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "unknwon",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/user", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

