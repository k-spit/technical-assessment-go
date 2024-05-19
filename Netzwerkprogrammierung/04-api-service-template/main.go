// Package main provides a simple CRUD API example using Gorilla Mux and PostgreSQL.
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"golang.org/x/time/rate"
)

// User represents a user in the system.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	db               *sql.DB
	cache            = Cache{items: make(map[string]CacheItem)}
	rateLimiter      = rate.NewLimiter(1, 3) // 1 request per second, burst size of 3
	cacheEnabled     bool
	rateLimitEnabled bool
)

// Cache represents a simple in-memory cache.
type Cache struct {
	mu    sync.Mutex
	items map[string]CacheItem
}

// CacheItem represents a single item in the cache.
type CacheItem struct {
	value      string
	expiration int64
}

// Set adds an item to the cache with an expiration duration.
func (c *Cache) Set(key, value string, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(duration).Unix()
	c.items[key] = CacheItem{
		value:      value,
		expiration: expiration,
	}
}

// Get retrieves an item from the cache.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found || item.expiration < time.Now().Unix() {
		return "", false
	}
	return item.value, true
}

func main() {
	apiURL := os.Getenv("API_URL")
	apiPort := os.Getenv("API_PORT")
	cacheEnabled, _ = strconv.ParseBool(os.Getenv("ENABLE_CACHE"))
	rateLimitEnabled, _ = strconv.ParseBool(os.Getenv("ENABLE_RATE_LIMITING"))

	var err error
	db, err = connectWithRetry(5, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/users", getUsersV1).Methods("GET")
	v1.HandleFunc("/users", createUser).Methods("POST")
	v1.HandleFunc("/users/{id}", getUser).Methods("GET")
	v1.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	v1.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	if cacheEnabled {
		v1.Use(cacheMiddleware)
	}
	if rateLimitEnabled {
		v1.Use(rateLimiterMiddleware)
	}

	serverAddress := fmt.Sprintf("%s:%s", apiURL, apiPort)
	fmt.Printf("Starting server on http://%s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}

// connectWithRetry attempts to connect to the database with retry logic.
func connectWithRetry(attempts int, sleep time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error
	dsn := os.Getenv("DATABASE_URL")
	for i := 0; i < attempts; i++ {
		db, err = sql.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to the database.")
				return db, nil
			}
		}
		log.Printf("Failed to connect to the database. Attempt %d/%d. Retrying in %v...", i+1, attempts, sleep)
		time.Sleep(sleep)
	}
	return nil, fmt.Errorf("could not connect to the database after %d attempts: %w", attempts, err)
}

// rateLimiterMiddleware is a middleware that limits the rate of requests.
func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// cacheMiddleware is a middleware that caches responses.
func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value, found := cache.Get(r.RequestURI); found {
			w.Write([]byte(value))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// getUsersV1 handles the GET /v1/users endpoint.
func getUsersV1(w http.ResponseWriter, r *http.Request) {
	if cacheEnabled {
		if value, found := cache.Get(r.RequestURI); found {
			w.Write([]byte(value))
			return
		}
	}

	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if users == nil {
		users = []User{} // Ensure we return an empty array instead of nil
	}

	response, _ := json.Marshal(users)
	if cacheEnabled {
		cache.Set(r.RequestURI, string(response), 10*time.Second)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// createUser handles the POST /v1/users endpoint.
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.QueryRow("INSERT INTO users (name) VALUES ($1) RETURNING id", user.Name).Scan(&user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// getUser handles the GET /v1/users/{id} endpoint.
func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if cacheEnabled {
		if value, found := cache.Get(r.RequestURI); found {
			w.Write([]byte(value))
			return
		}
	}

	var user User
	err = db.QueryRow("SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(user)
	if cacheEnabled {
		cache.Set(r.RequestURI, string(response), 10*time.Second)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// updateUser handles the PUT /v1/users/{id} endpoint.
func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("UPDATE users SET name = $1 WHERE id = $2", updatedUser.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	updatedUser.ID = id
	response, _ := json.Marshal(updatedUser)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// deleteUser handles the DELETE /v1/users/{id} endpoint.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
