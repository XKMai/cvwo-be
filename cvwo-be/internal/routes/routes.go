package routes

import (
	"net/http"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/handlers/comments"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/handlers/posts"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/handlers/users"
	"github.com/go-chi/chi/v5"
)

//Public Routes
func HealthCheckRoute() http.Handler {
	r := chi.NewRouter()

	// Define the health check endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	return r
}

func LoginRoute() chi.Router{
	r := chi.NewRouter()
	userhandler := users.UserHandler{}
	r.Post("/login", userhandler.LoginUser)
	r.Post("/create", userhandler.CreateUser)
	return r
}

//Protected Routes

func UserRoutes() chi.Router {
	r := chi.NewRouter()
	userhandler := users.UserHandler{}
	r.Get("/", userhandler.ListUsers) //Lists all
	r.Get("/{ID}", userhandler.GetUser)
	r.Put("/{ID}", userhandler.UpdateUser)
	r.Delete("/{ID}", userhandler.DeleteUser)
	r.Post("/refresh-token/", userhandler.RefreshToken)
	return r
}

func PostRoutes() chi.Router {
	r := chi.NewRouter()
	posthandler := posts.PostHandler{}
	r.Get("/", posthandler.ListPosts) //Pagination
	r.Post("/", posthandler.CreatePost)
	r.Get("/{ID}", posthandler.GetPost)
	r.Put("/{ID}", posthandler.UpdatePost)
	r.Delete("/{ID}", posthandler.DeletePost)
	return r
}

func CommentRoutes() chi.Router {
	r := chi.NewRouter()
	commenthandler := comments.CommentHandler{}
	r.Get("/", commenthandler.ListComments) //Pagination
	r.Post("/", commenthandler.CreateComment)
	r.Get("/{ID}", commenthandler.GetComment)
	r.Put("/{ID}", commenthandler.UpdateComment)
	r.Delete("/{ID}", commenthandler.DeleteComment)
	return r
}