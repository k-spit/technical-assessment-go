package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var mock sqlmock.Sqlmock

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/users", getUsersV1).Methods("GET")
	v1.HandleFunc("/users", createUser).Methods("POST")
	v1.HandleFunc("/users/{id}", getUser).Methods("GET")
	v1.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	v1.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	return r
}

func TestMain(m *testing.M) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestGetUsers(t *testing.T) {
	// Setup
	req, err := http.NewRequest("GET", "/v1/users", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router := setupRouter()

	// Mock DB response
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "John Doe").
		AddRow(2, "Jane Doe")
	mock.ExpectQuery("SELECT id, name FROM users").WillReturnRows(rows)

	// Execute
	router.ServeHTTP(rr, req)

	// Validate
	assert.Equal(t, http.StatusOK, rr.Code)
	var users []User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "Jane Doe", users[1].Name)
}

func TestCreateUser(t *testing.T) {
	// Setup
	user := User{Name: "John Doe"}
	userJSON, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router := setupRouter()

	// Mock DB response
	mock.ExpectQuery("INSERT INTO users \\(name\\) VALUES \\(\\$1\\) RETURNING id").
		WithArgs(user.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Execute
	router.ServeHTTP(rr, req)

	// Validate
	assert.Equal(t, http.StatusOK, rr.Code)
	var createdUser User
	err = json.Unmarshal(rr.Body.Bytes(), &createdUser)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", createdUser.Name)
}

func TestGetUser(t *testing.T) {
	// Setup
	req, err := http.NewRequest("GET", "/v1/users/1", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router := setupRouter()

	// Mock DB response
	row := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "John Doe")
	mock.ExpectQuery("SELECT id, name FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(row)

	// Execute
	router.ServeHTTP(rr, req)

	// Validate
	assert.Equal(t, http.StatusOK, rr.Code)
	var user User
	err = json.Unmarshal(rr.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", user.Name)
}

func TestUpdateUser(t *testing.T) {
	// Setup
	updatedUser := User{Name: "John Smith"}
	userJSON, _ := json.Marshal(updatedUser)
	req, err := http.NewRequest("PUT", "/v1/users/1", bytes.NewBuffer(userJSON))
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router := setupRouter()

	// Mock DB response
	mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
		WithArgs(updatedUser.Name, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	router.ServeHTTP(rr, req)

	// Validate
	assert.Equal(t, http.StatusOK, rr.Code)
	var user User
	err = json.Unmarshal(rr.Body.Bytes(), &user)
	assert.NoError(t, err)
	assert.Equal(t, "John Smith", user.Name)
}

func TestDeleteUser(t *testing.T) {
	// Setup
	req, err := http.NewRequest("DELETE", "/v1/users/1", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	router := setupRouter()

	// Mock DB response
	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute
	router.ServeHTTP(rr, req)

	// Validate
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
