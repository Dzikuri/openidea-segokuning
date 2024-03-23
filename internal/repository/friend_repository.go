package repository

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/pkg/errors"
)

type FriendRepository struct {
	DB *sql.DB
}

type RepositoryFriend interface {
	AddFriend(ctx context.Context, request model.FriendRequest) (*model.FriendResponse, error)
	CheckAlreadyFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, int, error)
	RemoveFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error)
}

func NewFriendRepository(db *sql.DB) RepositoryFriend {
	return &FriendRepository{
		DB: db,
	}
}

func (f *FriendRepository) AddFriend(ctx context.Context, request model.FriendRequest) (*model.FriendResponse, error) {

	dateTime := time.Now()
	dateCreate := dateTime.Format(time.RFC3339)
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Start a transaction with context
	tx, err := f.DB.BeginTx(context, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			return
		}
	}()

	// Insert into friends table
	queryCreate := `INSERT INTO friends (user_id, follow_user_id, created_at, updated_at) VALUES ($1, $2, $3, $4), ($2, $1, $3, $4) RETURNING id`
	_, err = tx.ExecContext(context, queryCreate, request.UserId, request.FriendId, dateCreate, dateCreate)
	if err != nil {
		return nil, err
	}

	// Update total_friend count in users table
	queryUpdateFollowerCount := `UPDATE users SET total_friend = total_friend + 1 WHERE id IN ($1, $2)`
	_, err = tx.ExecContext(context, queryUpdateFollowerCount, request.UserId, request.FriendId)
	if err != nil {
		return nil, err
	}

	// Commit the transaction if everything succeeded
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// If everything succeeded, return nil error
	return nil, nil
}

func (f *FriendRepository) RemoveFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, error) {

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Start a transaction with context
	tx, err := f.DB.BeginTx(context, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			return
		}
	}()

	// Delete from friends table
	queryDelete := `DELETE from friends where (user_id = $2 and follow_user_id = $1) or (user_id = $1 and follow_user_id = $2)`

	_, err = tx.ExecContext(context, queryDelete, userID, friendID)
	if err != nil {
		return nil, err
	}

	// Update total_friend count in users table
	queryUpdateFollowerCount := `UPDATE users SET total_friend = total_friend - 1 WHERE (id = $1 or id = $2)`
	_, err = tx.ExecContext(context, queryUpdateFollowerCount, userID, friendID)
	if err != nil {
		return nil, err
	}

	// Commit the transaction if everything succeeded
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (f *FriendRepository) CheckAlreadyFriend(ctx context.Context, userID string, friendID string) (*model.FriendResponse, int, error) {
	var count int

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var id string
	row := f.DB.QueryRowContext(context, "SELECT id FROM users WHERE id = $1", friendID)

	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, model.ErrResNotFound.Error
		}
		if strings.Contains(err.Error(), "invalid input syntax for type uuid") {
			// Handle specific error for invalid UUID format
			return nil, http.StatusNotFound, model.ErrUserNotFound
		}
		return nil, http.StatusInternalServerError, errors.Wrap(model.ErrInternalDatabase, err.Error())
	}

	row = f.DB.QueryRowContext(context, "SELECT count(*) FROM friends WHERE user_id = $1 AND follow_user_id = $2", userID, friendID)

	err = row.Scan(&count)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if count > 0 {
		return nil, 1, nil
	}

	return nil, 0, nil
}
