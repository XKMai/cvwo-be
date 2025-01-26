package main

import (
	"fmt"
	"net/http"
	"os"
	
	"github.com/joho/godotenv"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/router"
)

func main() {
	fmt.Println("Preparing server for connection...")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found. Falling back to system environment variables.")
	}

	db := database.SetupDatabase()
	r := router.Setup(db)
	env := os.Getenv("ENV")
	fmt.Println(env)
	if env == "DEV" {
		fmt.Println("Localhost hosted on port 3000")
		http.ListenAndServe(":3000", r)
		return
	} else if env == "PROD" {
		if err := http.ListenAndServeTLS(":3000", "server.crt", "server.key", r); err != nil {
		    fmt.Printf("Error starting server: %v\n", err)
		    os.Exit(1)
		}
	} else {
		fmt.Println("??")
	}

}
