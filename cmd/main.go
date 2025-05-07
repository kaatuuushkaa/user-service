package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/kaatuuushkaa/user-service/internal/auth"
	"github.com/kaatuuushkaa/user-service/internal/handler"
	"github.com/kaatuuushkaa/user-service/internal/storage"
	"log"
	"net/http"
	"time"
)

var (
	db  *sql.DB
	err error
)

func main() {
	for i := 0; i < 10; i++ {
		db, err = storage.NewPostgresDB()
		if err == nil {
			break
		}
		log.Println("Waiting for database...")
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	r := mux.NewRouter()

	public := r.PathPrefix("/").Subrouter()
	public.HandleFunc("/users/register", handler.RegisterUserHandler(db)).Methods("POST")

	private := r.PathPrefix("/").Subrouter()
	private.Use(auth.JWTMiddleware)

	private.HandleFunc("/users/{id}/status", handler.GetUserStatus(db)).Methods("GET")
	private.HandleFunc("/users/leaderboard", handler.GetLeaderboardHandler(db)).Methods("GET")
	private.HandleFunc("/users/{id}/task/complete", handler.CompleteTaskHandler(db)).Methods("POST")
	private.HandleFunc("/users/{id}/referrer", handler.AddReferrerHandler(db)).Methods("POST")

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	//time.Sleep(3600 * time.Second)
	log.Println("Server stopped.")
}
