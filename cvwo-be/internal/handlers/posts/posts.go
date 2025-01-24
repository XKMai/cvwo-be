package posts

// "encoding/json"
// "net/http"

// "github.com/XKMai/CVWO-React/CVWO-Backend/internal/database"

type PostHandler struct {}

// func (b *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	
// 	db := database.GetDB()
// 	var users []database.Post
// 	db.Find(&users)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(users)
// }

// func (b *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
// 	db := database.GetDB()
// 	id := r.URL.Query().Get("id")
// 	var user database.Post
// 	db.First(&user, id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(user)
// }

// func (b *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
// 	db := database.GetDB()
// 	var user database.Post
// 	json.NewDecoder(r.Body).Decode(&user)
// 	db.Create(&user)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(user)
// }

// func (b *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
// 	db := database.GetDB()
// 	var user database.Post
// 	json.NewDecoder(r.Body).Decode(&user)
// 	db.Save(&user)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(user)
// }

// func (b *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
// 	db := database.GetDB()
// 	id := r.URL.Query().Get("id")
// 	db.Delete(&database.Post{}, id)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// }
