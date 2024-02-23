package repository

import (
	"context"
	"database/sql"
	"net/http"
	"runtime/debug"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
)

type UserRepository interface {
	FindUserById(ctx context.Context, id int64) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, body *entity.User) (*entity.User, error)
	FindUserPassword(ctx context.Context, body *entity.User) (string, error)
	UpdateUserPassword(ctx context.Context, body *entity.User, password string) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}
func (r *userRepository) CreateUser(ctx context.Context, body *entity.User) (*entity.User, error) {
	user := entity.User{}

	hashedPassword, err := utils.HashPassword(body.Password, 12)
	if err != nil {
		return nil, err
	}
	runner := database.PickQuerier(ctx, r.db)
	q := `INSERT INTO users (email, name, password) 
	VALUES ($1, $2, $3 )
	returning id, email, name, created_at, updated_at, deleted_at`
	err = runner.QueryRowContext(ctx, q, body.Email, body.Name, hashedPassword).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &user, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	runner := database.PickQuerier(ctx, r.db)
	q := `SELECT u.id, u.email, u.name, u.created_at, u.updated_at from users u where u.email = $1 `

	err := runner.QueryRowContext(ctx, q, email).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, err.Error(), debug.Stack())
	}
	return &user, nil
}

func (r *userRepository) FindUserById(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	runner := database.PickQuerier(ctx, r.db)
	q := `SELECT u.id, u.email, u.name, u.created_at, u.updated_at from users u where u.id = $1 `

	err := runner.QueryRowContext(ctx, q, id).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, err.Error(), debug.Stack())
	}
	return &user, nil
}
func (r *userRepository) FindUserPassword(ctx context.Context, body *entity.User) (string, error) {

	q := `SELECT password from users where email = $1 `

	row := r.db.QueryRowContext(ctx, q, body.Email)
	var password string
	err := row.Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return password, nil
		}
		return "", apperror.NewInternalErrorType(http.StatusUnauthorized, constant.ResponseMsgErrorCredentials, debug.Stack())
	}
	return password, nil
}

func (r *userRepository) UpdateUserPassword(ctx context.Context, body *entity.User, password string) (*entity.User, error) {
	var user entity.User
	q := `UPDATE users SET password = $1, updated_at = now() where email = $2 RETURNING id, email, name, created_at, updated_at`
	hashedPassword, err := utils.HashPassword(password, 12)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, q, hashedPassword, body.Email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &user, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal, debug.Stack())
	}
	return &user, nil
}
