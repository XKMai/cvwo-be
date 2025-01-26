package posts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// "encoding/json"
// "net/http"

// "github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"

type PostHandler struct {}

func PaginateAndFilter(r *http.Request) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {

		type Input struct{
			Category string
			Page int
		}
		var input Input
		json.NewDecoder(r.Body).Decode(&input)

        // Add category filtering logic for pq.StringArray field
		category := input.Category
        if category != "" {
            db = db.Where("? = ANY(category)", category)
        }

        // Pagination logic
		page := input.Page
        if page <= 0 {
            page = 1
        }

        pageSize := 5 // Fixed page size of 5
        offset := (page - 1) * pageSize

        // Limit to 5 most recent items, ordered by a timestamp column (e.g., "created_at")
        return db.Preload("User",func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("password")
		}).Order("created_at DESC").Offset(offset).Limit(pageSize)
    }
}

//ListPosts does not get comments
func (b *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	var posts []models.Post
	if err := db.Scopes(PaginateAndFilter(r)).Find(&posts).Error; err != nil {
		http.Error(w, "Failed to find post: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}else{	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(posts)
	}
}

//GetPost will get comments as well
func (b *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var post models.Post
	post.ID = uint(id)
	if err := db.Preload("User",func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("password")}).Preload("Comments").First(&post).Error; err != nil {
		http.Error(w, "Failed to find post: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {
		json.NewEncoder(w).Encode(post)}
}

func (b *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	// Define the input structure
	type Input struct {
		Title string
		Content string
		Category string
		Picture string
		UserID uint
	}

	var input Input
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var categories_str = strings.Split(input.Category," ")
	// var categories, err = json.Marshal(categories_str)
	// if err != nil {
	// 	log.Fatalf("Failed to marshal category: %v", err)
	// }
	


	var user models.User
	if err := db.First(&user).Omit("password").Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} 

	// Create a new post
	newPost := models.Post{
		Title:   input.Title,
		Content: input.Content,
		Category: pq.StringArray(categories_str),
		UserID:  0,
		User:    user,     
		Score:   0,
		Picture: "", 
		Comments:[]models.Comment{},
	}

	// Save the post to the database
	if err := db.Create(&newPost).Error; err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} 

	// Respond with the created post
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newPost)
}

func (b *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	//Get ID from url param
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var post models.Post;
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil{		
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		} 

		//Convert to proper uint form
		post.ID = uint(id)

		// Save the post to the database
		if err := db.Save(&post).Preload("User",func(tx *gorm.DB) *gorm.DB {
			return tx.Omit("password")
		}).Error; err != nil {
			http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(nil)
			return
		}
	
		// Respond with the created post
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
}

func (b *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var post models.Post
	post.ID = uint(id)

	if err := db.Select("ID").First(&post).Error; err != nil {
		http.Error(w, "Failed to find post: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}

	if err := db.Delete(&post).Preload("User",func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("password")
	}).Error; err != nil {
		http.Error(w, "Failed to find post: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {json.NewEncoder(w).Encode(post)}
}