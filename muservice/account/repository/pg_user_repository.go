package repository

import (
	"context"
	"log"

	"github.com/NetworkPy/muserv/muservice/account/models"
	"github.com/NetworkPy/muserv/muservice/account/models/apperrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PGUserRepository is data/repository implementation
// of service layer UserRepository
type pgUserRepository struct {
	Db *sqlx.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *sqlx.DB) models.UserRepository {
	return &pgUserRepository{
		Db: db,
	}
}

// Create reaches out to database SQLX api
func (r *pgUserRepository) Create(ctx context.Context, u *models.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"

	if err := r.Db.GetContext(ctx, u, query, u.Email, u.Passowrd); err != nil {
		// check unique constraint
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

// FindByID fetches user by id
func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := "SELECT * FROM users WHERE uid=$1"

	// we need to actually check errors as it could be something other than not found
	if err := r.Db.GetContext(ctx, user, query, uid); err != nil {
		return user, apperrors.NewNotFound("uid", uid.String())
	}

	return user, nil
}
