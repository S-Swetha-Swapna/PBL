package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Database connection details
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "tripura88"
	DB_NAME     = "phonedb"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

var db *sql.DB

// Struct for incoming JSON request
type MobileNumber struct {
	Mobile string `json:"mobile"`
}

// Connect to PostgreSQL
func initDB()(*sql.DB, error) {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database:", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is not reachable:", err)
	}
	fmt.Println("Connected to PostgreSQL database successfully!")
	return db, nil
}

// Insert number into database
func insertMobileNumber(mobile string) error {
	_, err := db.Exec("INSERT INTO blocked_numbers (mobile) VALUES ($1)", mobile)
	return err
}

func jsonResponse(w http.ResponseWriter, message string, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}
// Handle POST request to store mobile number
func submitHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusOK)
        return
    }

    if r.Method != http.MethodPost {
       jsonResponse(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var data MobileNumber
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        jsonResponse(w, "Invalid input", http.StatusBadRequest)
        return
    }

    err = insertMobileNumber(data.Mobile)
    if err != nil {
        jsonResponse(w, "Failed to store number", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Number blocked successfully!"})
}


func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/submit", submitHandler)
	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
