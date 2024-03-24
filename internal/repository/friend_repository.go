package repository

import (
	"context"
	"database/sql"
	"fmt"
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
	FindAllFriend(ctx context.Context, request model.GetFriendListRequest) ([]model.FriendResponse, model.MetaDataResponse, error)
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

func (f *FriendRepository) FindAllFriend(ctx context.Context, request model.GetFriendListRequest) ([]model.FriendResponse, model.MetaDataResponse, error) {

	friends := make([]model.FriendResponse, 0)

	queryCondition := fmt.Sprintf(" WHERE users.id <> '%s' ", request.UserId)
	queryJoin := ""
	if request.OnlyFriend == true {
		queryCondition += fmt.Sprintf(" AND friends.user_id = '%s' ", request.UserId)

		// NOTE Using Join
		queryJoin += " INNER JOIN friends ON users.id = friends.follow_user_id "
	}

	if request.Search != "" {
		queryCondition += fmt.Sprintf(" AND users.name LIKE '%%%s%%' ", request.Search)
	}

	orderBy := "DESC"
	if request.OrderBy != "" {
		orderBy = request.OrderBy
	}

	sortBy := "created_at"
	if request.SortBy != "" {
		if request.SortBy == "friendCount" {
			sortBy = "total_friend"
		} else {
			sortBy = "created_at"
		}
	}

	querySortBy := fmt.Sprintf(" ORDER BY %s %s ", sortBy, orderBy)

	// NOTE Using Join
	queryGet := fmt.Sprintf("SELECT users.id, users.name, users.image_url, users.total_friend, users.created_at FROM users %s %s %s LIMIT $1 OFFSET $2", queryJoin, queryCondition, querySortBy)
	fmt.Println(queryGet)
	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	rows, err := f.DB.QueryContext(context, queryGet, request.Limit, request.Offset)
	if err != nil {
		return nil, model.MetaDataResponse{}, err
	}

	totalRows := 0
	for rows.Next() {
		var friend model.FriendResponse

		var createdAt time.Time

		err = rows.Scan(&friend.UserId, &friend.Name, &friend.ImageUrl, &friend.FriendCount, &createdAt)
		if err != nil {
			return nil, model.MetaDataResponse{}, err
		}

		// friend.CreatedAt = createdAt.Format("")
		friend.CreatedAt, _ = time.Parse(time.RFC3339, friend.CreatedAt.String())
		friends = append(friends, friend)

		totalRows++
	}
	defer rows.Close()

	return friends, model.MetaDataResponse{
		Limit:  request.Limit,
		Offset: request.Offset,
		Total:  totalRows,
	}, nil
}
