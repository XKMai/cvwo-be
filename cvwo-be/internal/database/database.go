package database

import (
	"fmt"
	"log"

	//"github.com/XKMai/CVWO-React/CVWO-Backend/internal/handlers/users"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func SetupDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=Abc123!@# dbname=cvwo_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	fmt.Println("Database connection established!")	
	db.AutoMigrate(&models.User{},&models.Post{},&models.Comment{})
	return db
}