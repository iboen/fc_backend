package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	// Open a connection to the MySQL database
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET")
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	// Query to retrieve the latest image URL
	query := "SELECT url FROM images ORDER BY created_at DESC LIMIT 1"

	// Execute the query
	var imageURL string
	err := db.QueryRow(query).Scan(&imageURL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatalf("Error retrieving latest image URL: %v", err)
		return
	}

	// Enable CORS
	enableCors(&w)

	// Set the Content-Type header to indicate that the response contains JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response containing the image URL to the response writer
	w.Write([]byte(`{"data": "` + imageURL + `"}`))

	// Log the request
	log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response status code to indicate that the application is healthy (200 OK)
	w.WriteHeader(http.StatusOK)
	// Write a simple message indicating that the application is healthy
	w.Write([]byte("OK"))
}

func main() {
	// Initialize the database connection
	initDB()
	defer db.Close()

	// Define the route and handler for the /get-image endpoint
	http.HandleFunc("/get-image", getImageHandler)

	// Define the route and handler for the health check endpoint (/)
	http.HandleFunc("/", healthCheckHandler)

	// get port from env, if not set default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the HTTP server
	log.Println("Server listening on port ", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}
