package repository

import (
	"context"
	"database/sql"
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/entity"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/constant"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/database"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/utils"
)

type PasswordTokenRepository interface {
	// FindAll(ctx context.Context) ([]entity.User, error)
	// FindSimilarUserByName(ctx context.Context, name string) ([]entity.User, error)
	// FindUserById(ctx context.Context, id int64) (*entity.User, error)
	CheckToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error)
	CreateToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error)
	UpdateDeleteToken(ctx context.Context, body *entity.User, token string) (*entity.PasswordToken, error)
	GetValidToken(ctx context.Context, token string) (*entity.PasswordToken, error)
}

type passwordTokenRepository struct {
	db *sql.DB
}

func NewPasswordTokenRepository(db *sql.DB) *passwordTokenRepository {
	return &passwordTokenRepository{
		db: db,
	}
}
func (r *passwordTokenRepository) CreateToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error) {
	passToken := entity.PasswordToken{}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	runner := database.PickQuerier(ctx, r.db)

	q := `INSERT INTO password_tokens (user_id, token) 
	VALUES ($1, $2)
	returning id, user_id, token, expired_at, created_at, updated_at, deleted_at`

	err = runner.QueryRowContext(ctx, q, body.ID, token).Scan(&passToken.ID, &passToken.UserID, &passToken.Token, &passToken.ExpiredAt, &passToken.CreatedAt, &passToken.UpdatedAt, &passToken.DeletedAt)
	if err != nil {
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}

	return &passToken, nil
}

func (r *passwordTokenRepository) CheckToken(ctx context.Context, body *entity.User) (*entity.PasswordToken, error) {
	passToken := entity.PasswordToken{}
	runner := database.PickQuerier(ctx, r.db)
	q := `SELECT id, user_id, token, expired_at, created_at, updated_at, deleted_at from password_tokens where user_id = $1 and deleted_at is null`

	err := runner.QueryRowContext(ctx, q, body.ID).Scan(&passToken.ID, &passToken.UserID, &passToken.Token, &passToken.ExpiredAt, &passToken.CreatedAt, &passToken.UpdatedAt, &passToken.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &passToken, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	return &passToken, nil
}

func (r *passwordTokenRepository) UpdateDeleteToken(ctx context.Context, body *entity.User, token string) (*entity.PasswordToken, error) {
	passToken := entity.PasswordToken{}
	runner := database.PickQuerier(ctx, r.db)
	q := `UPDATE password_tokens SET deleted_at = now() where user_id = $1 and token = $2 returning id, user_id, token, expired_at, created_at, updated_at, deleted_at`

	err := runner.QueryRowContext(ctx, q, body.ID, token).Scan(&passToken.ID, &passToken.UserID, &passToken.Token, &passToken.ExpiredAt, &passToken.CreatedAt, &passToken.UpdatedAt, &passToken.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &passToken, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	return &passToken, nil
}

func (r *passwordTokenRepository) GetValidToken(ctx context.Context, token string) (*entity.PasswordToken, error) {
	passToken := entity.PasswordToken{}
	runner := database.PickQuerier(ctx, r.db)
	q := `SELECT id, user_id, token, expired_at, created_at, updated_at, deleted_at from password_tokens where token = $1 and expired_at < now()`

	err := runner.QueryRowContext(ctx, q, token).Scan(&passToken.ID, &passToken.UserID, &passToken.Token, &passToken.ExpiredAt, &passToken.CreatedAt, &passToken.UpdatedAt, &passToken.DeletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &passToken, nil
		}
		return nil, apperror.NewInternalErrorType(http.StatusInternalServerError, constant.ResponseMsgErrorInternal)
	}
	return &passToken, nil
}
