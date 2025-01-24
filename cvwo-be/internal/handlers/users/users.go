package users

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"gorm.io/gorm"
)

type UserHandler struct {
}



func (b *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	fmt.Println("ListUsers")
	var users []models.User
	// Get all records
	result := db.Select("Name").Find(&users)
	// SELECT * FROM users;
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (b *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers")
	// db := database.GetDB()
	// id := r.URL.Query().Get("id")
	// var user database.User
	// db.First(&user, id)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(user)
}

func (b *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	// Define the input structure
	type Input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var input Input
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Create a new user
	newUser := models.User{
		Role:        "User",
		Name:        input.Name,
		Password:    input.Password,
		Picture:     byte(0),
		Description: "",
	}

	// Save the user to the database
	if err := db.Create(&newUser).Error; err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

func (b *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateUser")
	// db := database.GetDB()
	// var user database.User
	// json.NewDecoder(r.Body).Decode(&user)
	// db.Save(&user)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(user)
}

func (b *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteUsers")
	// db := database.GetDB()
	// id := r.URL.Query().Get("id")
	// db.Delete(&database.User{}, id)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
}

