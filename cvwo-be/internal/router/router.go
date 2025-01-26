package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/contextkeys"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/routes"
	"github.com/XKMai/CVWO-React/CVWO-Backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

// Setup initializes and returns the router with all routes and middleware
func Setup(db *gorm.DB) chi.Router {
	r := chi.NewRouter()
	fmt.Println(db)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(InjectDB(db)) // Inject the database into the request context

	apiRouter := chi.NewRouter()
	apiRouter.Use(InjectDB(db))
	apiRouter.Use(CheckDBInContext) // Check if db is set in context

	// CORS setup
	apiRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},                     // List of allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},  // HTTP methods allowed
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},  // Headers allowed in requests
		ExposedHeaders:   []string{"Link"},                                     // Headers exposed to clients
		AllowCredentials: true,                                                // Allow cookies and credentials
		MaxAge:           300,                                                 // Cache duration for preflight requests
	}))

	// Public routes (e.g., health check, login)
	apiRouter.Mount("/healthcheck", routes.HealthCheckRoute())
	apiRouter.Mount("/users", routes.LoginRoute())

	// Protected routes
	protectedRouter := chi.NewRouter()
	protectedRouter.Use(utils.IsAuthorized) // Apply the authorization middleware
	protectedRouter.Use(InjectDB(db))

	protectedRouter.Mount("/users", routes.UserRoutes())        // Protected user routes
	protectedRouter.Mount("/posts", routes.PostRoutes())        // Protected post routes
	protectedRouter.Mount("/comments", routes.CommentRoutes())  // Protected comment routes

	// Mount the protected routes under /api/protected
	apiRouter.Mount("/protected", protectedRouter)

	// Mount the API router to the main router
	r.Mount("/api", apiRouter)

	return r
}

// InjectDB middleware adds the database instance to the request context
func InjectDB(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextkeys.DBContextKey, db) // Add db to context
			next.ServeHTTP(w, r.WithContext(ctx))                  // Pass the updated context
		})
	}
}

func CheckDBInContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the DB from the context
		db, ok := r.Context().Value(contextkeys.DBContextKey).(*gorm.DB)
		if !ok || db == nil {
			// If the DB is not set correctly, log an error or return a response indicating the issue
			fmt.Println("ERROR: Database is not set correctly in context.")
			http.Error(w, "Database is not set correctly in context.", http.StatusInternalServerError)
			return
		}

		// If DB is set, continue with the request
		fmt.Println("Database is correctly set in context.")
		next.ServeHTTP(w, r)
	})
}
