package main

import (
	//"net/http"

	//"fmt"

	"net/http"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"
	//"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/router"
	//"github.com/XKMai/CVWO-React/CVWO-Backend/internal/handlers/users"
	//"gorm.io/driver/postgres"
	//"gorm.io/gorm"
)

func main() {
	db := database.SetupDatabase()
	r := router.Setup(db)
	//db, err := gorm.Open(postgres.Open("cvwo_db"), &gorm.Config{})
	//user := models.User{ID: 1, Role:"User", Name: "Xin Kai", Password: "password"}
	http.ListenAndServe(":3000", r)
}
