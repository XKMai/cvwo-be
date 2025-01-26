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
	env := os.Getenv("ENV")
	if env == "DEV" {
		http.ListenAndServe(":3000", r)
		return
	} else if env == "PROD" {
		http.ListenAndServeTLS(":3000", "cert.pem", "key.pem", r)
	}
}
