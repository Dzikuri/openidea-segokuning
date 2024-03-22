package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/Dzikuri/openidea-segokuning/internal/helper"
	"github.com/Dzikuri/openidea-segokuning/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type RepositoryUser interface {
	Register(ctx context.Context, user *model.UserAuthRequest) (*model.UserResponse, error)
	FindByPhone(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error)
	FindByEmail(ctx context.Context, user *model.UserAuthRequest) (exists bool, res *model.UserResponse, err error)
	FindById(ctx context.Context, id string) (res *model.UserResponse, code int, err error)
	UpdateUserData(ctx context.Context, request model.UserResponse) (*model.UserResponse, error)
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
	helper.LogPretty(user)
	querySelect := fmt.Sprintf("SELECT id, email, phone, name, password, created_at, updated_at FROM users WHERE phone = $1")

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result := &model.UserResponse{}
	row := r.DB.QueryRowContext(context, querySelect, user.CredentialValue)

	errRowScan := row.Scan(&result.Id, &result.Email, &result.Phone, &result.Name, &result.Password, &result.CreatedAt, &result.UpdatedAt)
	fmt.Println("err get data from database : ", errRowScan)
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
	fmt.Println("err get data from database : ", errRowScan)
	if errors.Is(errRowScan, sql.ErrNoRows) {
		return false, nil, model.ErrUserNotFound
	}
	if errRowScan != nil {
		return false, nil, err
	}

	return true, result, nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (res *model.UserResponse, code int, err error) {

	var user model.UserResponse

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(context, "SELECT id, email, phone, name,  password, created_at, updated_at FROM users WHERE id = $1", id)

	err = row.Scan(&user.Id, &user.Email, &user.Phone, &user.Name, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, errors.Wrap(model.ErrUserNotFound, model.ErrUserNotFound.Error())
		}
		return nil, http.StatusInternalServerError, errors.Wrap(model.ErrInternalDatabase, err.Error())
	}

	return &user, http.StatusOK, nil
}

func (r *UserRepository) UpdateUserData(ctx context.Context, request model.UserResponse) (*model.UserResponse, error) {

	queryUpdate := "UPDATE users SET"
	var values []interface{}
	counter := 1

	if request.Email != "" {
		queryUpdate += fmt.Sprintf(" email = $%d,", counter)
		values = append(values, request.Email)
		counter++
	}

	if request.Phone != "" {
		queryUpdate += fmt.Sprintf(" phone = $%d,", counter)
		values = append(values, request.Phone)
		counter++
	}

	if request.Name != "" {
		queryUpdate += fmt.Sprintf(" name = $%d,", counter)
		values = append(values, request.Name)
		counter++
	}

	if request.ImageUrl != "" {
		queryUpdate += fmt.Sprintf(" image_url = $%d,", counter)
		values = append(values, request.ImageUrl)
		counter++
	}

	queryUpdate += fmt.Sprintf(" updated_at = $%d,", counter)
	values = append(values, time.Now())
	counter++

	// Remove the last comma if present
	if counter > 1 {
		queryUpdate = queryUpdate[:len(queryUpdate)-1]
	}

	queryUpdate += fmt.Sprintf(" WHERE id = $%d", counter)
	values = append(values, request.Id)

	context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(context, queryUpdate, values...)
	if err != nil {
		return nil, errors.Wrap(model.ErrInternalDatabase, err.Error())
	}

	row, err := result.RowsAffected()

	if err != nil {
		return nil, errors.Wrap(echo.ErrInternalServerError, err.Error())
	}

	if row == 0 {
		return nil, err
	}

	return nil, nil
}
