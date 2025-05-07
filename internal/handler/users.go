package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kaatuuushkaa/user-service/internal/auth"
	"github.com/kaatuuushkaa/user-service/internal/storage"
	"net/http"
	"strconv"
)

func GetUserStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
			return
		}

		user, err := storage.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func GetLeaderboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		leadboard, err := storage.GetLeaderboard(db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching leaderboard: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(leadboard)
	}
}

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		userID, err := storage.CreateUser(db, req.Name)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creatimg user: %v", err), http.StatusInternalServerError)
			return
		}

		token, err := auth.GenerateToken(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error generating token: %v", err), http.StatusInternalServerError)
			return
		}

		resp := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
