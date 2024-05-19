package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users   = make(map[int]*User)
	nextID  = 1
	muxLock sync.Mutex
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	muxLock.Lock()
	userList := make([]*User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	muxLock.Unlock()
	json.NewEncoder(w).Encode(userList)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	muxLock.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = &user
	muxLock.Unlock()
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	muxLock.Lock()
	user, exists := users[id]
	muxLock.Unlock()
	if !exists {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	muxLock.Lock()
	user, exists := users[id]
	if !exists {
		muxLock.Unlock()
		http.NotFound(w, r)
		return
	}
	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		muxLock.Unlock()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Name = updatedUser.Name
	muxLock.Unlock()
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	muxLock.Lock()
	_, exists := users[id]
	if exists {
		delete(users, id)
		muxLock.Unlock()
		w.WriteHeader(http.StatusNoContent)
	} else {
		muxLock.Unlock()
		http.NotFound(w, r)
	}
}
