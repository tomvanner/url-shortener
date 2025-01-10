package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db     *sql.DB
	config Config = Config{
		Env:         getEnv("APP_ENV", "DEV"),
		Port:        getEnv("APP_PORT", "8000"),
		DatabaseDSN: getEnv("DB_PATH", "database/urls.db"),
	}
)

type Config struct {
	Env         string
	Port        string
	DatabaseDSN string
}

type CreateData struct {
	LongURL string `json:"longURL"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", config.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Database connection established")
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var d CreateData
	err := json.NewDecoder(r.Body).Decode(&d)

	if err != nil || !isValidURL(d.LongURL) {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	slug := generateSlug()

	// @todo: check if slug already exists
	_, err = db.Exec(
		"INSERT INTO urls (long_url, slug, created) VALUES ($1, $2, $3)",
		d.LongURL,
		slug,
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"slug": slug,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	shortURL := make([]byte, length)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	slug := r.URL.Path[1:]

	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE slug = $1", slug).Scan(&longURL)

	if err == sql.ErrNoRows {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Fatal(err)
		http.Error(w, "There was an error. Please try again later.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func isValidURL(urlString string) bool {
	u, err := url.Parse(urlString)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}

	return value
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/create", createShortURLHandler)

	http.ListenAndServe(":"+config.Port, nil)
}
