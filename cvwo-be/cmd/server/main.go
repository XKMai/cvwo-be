package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/router"
)

func main() {
	fmt.Println("Preparing server for connection...")
	db := database.SetupDatabase()
	r := router.Setup(db)
	fmt.Println("Hehe")
	env := os.Getenv("ENV")
	if env == "DEV" {
		fmt.Println("Localhost hosted on port 3000")
		http.ListenAndServe(":3000", r)
		return
	} else if env == "PROD" {
		fmt.Println("Production server hosted on port 3000")
		http.ListenAndServeTLS(":443", "cert.pem", "key.pem", r)
		return
	}
}
