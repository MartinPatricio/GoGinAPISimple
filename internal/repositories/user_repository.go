package repositories

import (
	"context"
	"errors"
	"firstapi/internal/models"
	"firstapi/pkg/db"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

/*Definimo las operaciones en nuetro contrato*/
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	/*GetById(ctx context.Context, idUser string) (*models.User, error)
	GetByEmail(ctx context.Context, Email string) (*models.User, error)
	GetAll(ctx context.Context, filter *models.UserFilter, pagination *models.Pagination) ([]*models.User, int, error)
	Update(ctx context.Context, idUser string, user *models.User) error
	Delete(ctx context.Context, idUser string) error

	//Specific operations by Models
	UpdatePassword(ctx context.Context, idUser string, passHash string) error
	GetRoles(ctx context.Context) ([]*models.Rol, error)*/
}

//Implementing PostgreSql

type userRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
	INSERT INTO users (idUser, idRol, NameUser, Email, LastName, LastActivitie, DateCreated, Password)
        VALUES ($1, $2, $3, $4, $5, $6 ,$7, $8)
        RETURNING id_user, created_at, updated_at`
	err := ur.db.QueryRow(ctx, query,
		user.IdUser,
		user.IdRol,
		user.NameUser,
		user.Email,
		user.LastName,
		user.LastActivitie,
		user.DateCreated,
		user.PasswordHash,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			if pqErr.Code == "23505" {
				if strings.Contains(pqErr.Error(), "idUser") {
					return errors.New("Id User already exists")
				}
			}
		} else {
			return fmt.Errorf("Failed to create user: %w", err)
		}
	}
	return nil
}
