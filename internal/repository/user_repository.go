package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type RepositoryUser interface {
	Register(ctx context.Context, user *model.UserAuthRequest) (*model.UserResponse, error)
	FindByPhone(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error)
	FindByEmail(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error)
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) RepositoryUser {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Register(ctx context.Context, user *model.UserAuthRequest) (*model.UserResponse, error) {

	queryCreate := ""
	if user.CredentialType == model.Email {
		queryCreate = `
            INSERT INTO users(email, name, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id
        `
	}

	if user.CredentialType == model.Phone {
		queryCreate = `
            INSERT INTO users(phone, name, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id
        `
	}

	// queryCreate := fmt.Sprintf("INSERT INTO users(%s, name, password, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id", insertTo)

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := r.DB.QueryRowContext(context, queryCreate, &user.CredentialValue, &user.Name, &user.Password, time.Now(), time.Now())
	var id = ""
	err := result.Scan(&id)
	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, model.ErrUserAlreadyExists
			default:
				return nil, err
			}
		}
		return nil, err
	}

	res := new(model.UserResponse)
	if user.CredentialType == model.Email {
		res.Email = user.CredentialValue
		res.Phone = ""
	}

	if user.CredentialType == model.Phone {
		res.Phone = user.CredentialValue
		res.Email = ""
	}

	res.Id = helper.GetUUID(id)
	res.Name = user.Name

	return res, nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error) {
	querySelect := fmt.Sprintf("SELECT id, email, phone, name, password, created_at, updated_at FROM users WHERE phone = $1")

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := new(model.UserResponse)

	errRowScan := r.DB.QueryRowContext(context, querySelect, user.CredentialValue).Scan(&result.Id, &result.Email, &result.Phone, &result.Password, &result.CreatedAt, &result.UpdatedAt)

	if errors.Is(errRowScan, sql.ErrNoRows) {
		return false, nil, model.ErrUserNotFound
	}

	if errRowScan != nil {
		return false, nil, err
	}

	return true, result, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error) {
	querySelect := fmt.Sprintf("SELECT id, email, phone, name, password, created_at, updated_at FROM users WHERE email = $1")

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := &model.UserResponse{}
	row := r.DB.QueryRowContext(context, querySelect, user.CredentialValue)

	errRowScan := row.Scan(&result.Id, &result.Email, &result.Phone, &result.Name, &result.Password, &result.CreatedAt, &result.UpdatedAt)
	if errors.Is(errRowScan, sql.ErrNoRows) {
		return false, nil, model.ErrUserNotFound
	}
	if errRowScan != nil {
		return false, nil, err
	}

	return true, result, nil
}
