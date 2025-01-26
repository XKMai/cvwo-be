package contextkeys

// Define a custom type to avoid collisions
type contextKey string

// DBContextKey is the exported key for the database object
const DBContextKey contextKey = "db"