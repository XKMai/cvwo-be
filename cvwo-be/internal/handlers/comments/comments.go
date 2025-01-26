package comments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// "encoding/json"
// "net/http"

// "github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"

type CommentHandler struct {}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {

		type Input struct{
			Page int
			PostID int `json:"post_id"`
		}
		var input Input
		json.NewDecoder(r.Body).Decode(&input)

		var page = input.Page
        // Default values if not provided
        if page <= 0 {
            page = 1
        }
		
        pageSize := 5 // Fixed page size of 5
        offset := (page - 1) * pageSize

        // Query for pagination
        return db.
            Where("post_id = ?", input.PostID).   // Filter comments by PostID
            Order("created_at ASC").        // Oldest to newest
            Offset(offset).                 // Apply the offset
            Limit(pageSize).                // Limit results to the page size
            Preload("User",func(tx *gorm.DB) *gorm.DB {
				return tx.Omit("password")
			})                 // Preload the "User" relation (optional)
    }
}

func (b *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	var comments []models.Comment
	if err := db.Scopes(Paginate(r)).Find(&comments).Error; err != nil {
		http.Error(w, "Failed to find comment: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}else{	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comments)
	}
}

func (b *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var comment models.Comment
	comment.ID = uint(id)
	if err := db.Preload("User").First(&comment).Error; err != nil {
		http.Error(w, "Failed to find comment: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {
		comment.User.Password = ""
		json.NewEncoder(w).Encode(comment)}
}

func (b *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)

	// Define the input structure
	type Input struct {
		Content string `json:"content"`
		UserID  uint   `json:"user_id"`
		PostID  uint   `json:"post_id"`
	}

	var input Input
	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Fetch the user
	var user models.User
	if err := db.First(&user, "id = ?", input.UserID).Error; err != nil {
		http.Error(w, "Failed to find user: "+err.Error(), http.StatusNotFound)
		return
	}
	user.Password = ""

	// Fetch the post
	var post models.Post
	if err := db.First(&post, "id = ?", input.PostID).Error; err != nil {
		http.Error(w, "Failed to find post: "+err.Error(), http.StatusNotFound)
		return
	}

	// Create a new comment
	newComment := models.Comment{
		Content: input.Content,
		UserID:  user.ID, // Correctly assign the user ID
		PostID:  post.ID, // Correctly assign the post ID
		User:    user,    // Include the user model if needed
	}

	// Save the comment to the database
	if err := db.Model(&post).Association("Comments").Append(&newComment); err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	post.Comments = append(post.Comments, newComment)

	// Respond with the created comment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newComment)
}

func (b *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	//Get ID from url param
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var comment models.Comment;
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil{		
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		} 

		//Convert to proper uint form
		comment.ID = uint(id)

		// Save the comment to the database
		if err := db.Save(&comment).Error; err != nil {
			http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
			json.NewEncoder(w).Encode(nil)
			return
		}
	
		// Respond with the created comment
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(comment)
}

func (b *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("db").(*gorm.DB)
	id_str := chi.URLParam(r, "ID")
	id,err := strconv.ParseUint(id_str,10,64)

	if err != nil{
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
		}

	var comment models.Comment
	comment.ID = uint(id)

	if err := db.Select("ID").First(&comment).Error; err != nil {
		http.Error(w, "Failed to find comment: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	}

	if err := db.Delete(&comment).Error; err != nil {
		http.Error(w, "Failed to find comment: "+err.Error(), http.StatusInternalServerError)
		json.NewEncoder(w).Encode(nil)
		return
	} else {json.NewEncoder(w).Encode(comment)}
}