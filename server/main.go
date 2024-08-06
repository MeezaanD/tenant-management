package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"tenant-management/database"
	"tenant-management/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Define a template directly in the code
var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
</head>
<body>
    <h1>Welcome to Our Tenant Management System</h1>
    <p>This page is dynamically generated by Go's html/template package.</p>
</body>
</html>
`))

func main() {
	database.Init()

	router := mux.NewRouter()

	// Set up routes for API
	routes.SetupRoutes(router)

	// Handle root URL by serving a dynamic HTML page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start the server on the port specified in the .env file
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"  // Default to port 8080 if nothing is specified
	}
	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
