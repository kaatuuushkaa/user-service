package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kaatuuushkaa/user-service/internal/storage"
	"net/http"
	"strconv"
)

func CompleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
			return
		}

		err = storage.CompleteTask(db, userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error completing task: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func AddReferrerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid user ID: %v", err), http.StatusBadRequest)
			return
		}

		//чтение реф кода из тела запроса
		var req struct {
			ReferrerID int `json:"referrer_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
			return
		}

		err = storage.AddReferrer(db, userID, req.ReferrerID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding referre: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
