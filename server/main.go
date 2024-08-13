package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"tenant-management/database"
	"tenant-management/routes"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
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
	routes.SetupRoutes(router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	
	corsOptions := handlers.AllowedOrigins([]string{"*"}) // Allow all origins for development
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	log.Printf("Server started at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsOptions, corsMethods, corsHeaders)(router)))
}
