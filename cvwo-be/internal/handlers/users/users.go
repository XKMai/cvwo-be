package users

import (
	//"encoding/json"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/contextkeys"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type UserHandler struct {
}

func (b *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}

	var users []models.User
	if err := db.Omit("password").Find(&users).Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}else{	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	}
}

func (b *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var user models.User
	user.ID = uint(id)
	if err := db.Omit("password").First(&user).Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {
		json.NewEncoder(w).Encode(user)}
}

func (b *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}

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

	//Hashes Password
	hashedPassword,err := utils.HashPassword(input.Password)
	if err != nil{
		http.Error(w, "Error hashing password",http.StatusInternalServerError)
		return
	}

	// Create a new user
	newUser := models.User{
		Role:        "User",
		Name:        input.Name, 
		Password:    hashedPassword,
		Description: "",
	}

	if newUser.Name == "Admin"{newUser.Role = "Admin"} //Initialises admin for Admin

	// Save the user to the database
	if err := db.Create(&newUser).Omit("password").Error; err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} 

	// Respond with the created user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newUser)
}

//Only allow updating Description of User
func (b *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}
	//Get ID from url param
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	type Input struct {
		Description string `json:"description"`
	}
	var input Input
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil{		
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		} 

		var user models.User
		//Convert to proper uint form
		user.ID = uint(id)
		user.Description = input.Description

		// Save the user to the database
		if err := db.Save(&user).Omit("password").Error; err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(nil)
			return
		}
	
		// Respond with the created user
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
}

func (b *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var user models.User
	user.ID = uint(id)

	if err := db.Select("ID").First(&user).Omit("password").Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {json.NewEncoder(w).Encode(user)}
}

// LoginUser - Create JWT token on user login
func (b *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}

	type LoginInput struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var input LoginInput
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Find user by name
	var user models.User
	if err := db.Where("name = ?",input.Name).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Compare passwords (you need to hash the input password and check it)
	if err := utils.CheckPassword(input.Password, user.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Create JWT Access Token
	accessToken, err := utils.CreateToken(user.ID, user.Name)
	if err != nil {
		http.Error(w, "Error creating access token", http.StatusInternalServerError)
		return
	}

	// Create JWT Refresh Token
	refreshToken, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}

	// Respond with both the access and refresh tokens
	res := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		UserID int `json:"user_id"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID: int(user.ID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (b *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
	if !ok || db == nil {
		fmt.Println("DB is nil or not properly set in context")
		http.Error(w, "Database connection is not available", http.StatusInternalServerError)
		return
	}

	// Extract the refresh token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Missing Authorization token", http.StatusUnauthorized)
		return
	}

	// Remove 'Bearer ' prefix if it exists
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Verify the refresh token
	claims, err := utils.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Find the user based on the claims (user ID in this case)
	var user models.User
	if err := db.First(&user, claims.Subject).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate a new access token
	accessToken, err := utils.CreateToken(user.ID, user.Name)
	if err != nil {
		http.Error(w, "Error creating new token", http.StatusInternalServerError)
		return
	}

	// Generate a new refresh token (optional)
	refreshToken, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}

	// Respond with the new tokens
	res := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}