package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
)

var (
	urlStore = make(map[string]string)
	mu       sync.Mutex
)

type CreateData struct {
	LongURL string `json:"longURL"`
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var d CreateData
	err := json.NewDecoder(r.Body).Decode(&d)

	// @todo: add url validation
	if err != nil || d.LongURL == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()

	mu.Lock()
	urlStore[shortURL] = d.LongURL
	mu.Unlock()

	response := map[string]string{
		"short_url": "http://localhost:8000/" + shortURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateShortURL() string {
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

	shortURL := r.URL.Path[1:]

	mu.Lock()
	longURL, exists := urlStore[shortURL]
	mu.Unlock()

	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func main() {
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/create", createShortURLHandler)

	http.ListenAndServe(":8000", nil)
}
