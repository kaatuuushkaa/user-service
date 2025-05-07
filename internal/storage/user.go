package storage

//работа с данными пользователей

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kaatuuushkaa/user-service/internal/models"
	_ "github.com/lib/pq"
)

func CreateUser(db *sql.DB, name string) (int, error) {
	var userID int

	err := db.QueryRow(
		"INSERT INTO users(name) VALUES ($1) RETURNING id", name).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("Error creating user: %v", err)
	}

	return userID, nil
}

func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var user models.User

	err := db.QueryRow("SELECT id, name, points FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User with Id %d not found", id)
		}
		return nil, fmt.Errorf("Could not get user: %w", err)
	}

	return &user, nil
}

func GetLeaderboard(db *sql.DB) ([]*models.User, error) {
	rows, err := db.Query("SELECT id, name, points FROM users ORDER BY points DESC LIMIT 3")
	if err != nil {
		return nil, fmt.Errorf("Could not get leaderboard: %w", err)
	}
	defer rows.Close()

	var leaderboard []*models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Points); err != nil {
			return nil, fmt.Errorf("Could not scan user data: %w", err)
		}
		leaderboard = append(leaderboard, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while iterating over rows: %w", err)
	}

	return leaderboard, nil
}

func CompleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("UPDATE users SET points = points +1 WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Could not update user points: %w", err)
	}
	return nil
}

func AddReferrer(db *sql.DB, userID int, referrerID int) error {
	if userID == referrerID {
		return errors.New("User cannot refer themselves")
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", referrerID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Referrer user not found")
	}

	var currentReferrer sql.NullInt64
	err = db.QueryRow("SELECT referrer_id FROM users WHERE id = $1", userID).Scan(&currentReferrer)
	if err != nil {
		return err
	}
	if currentReferrer.Valid {
		return errors.New("Referrer already set for this user")
	}

	_, err = db.Exec("UPDATE users SET referrer_id = $1 WHERE id = $2", referrerID, userID)
	return err
}
