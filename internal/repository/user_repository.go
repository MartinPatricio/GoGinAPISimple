package repository

import (
	"context"
	"database/sql"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
)

// UserRepository define la interfaz para las operaciones de usuario.
type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Tbluser, error)
	DeleteUser(ctx context.Context, idUser int32) error
	GetUserByID(ctx context.Context, idUser int32) (db.Tbluser, error)
	GetUserByEmail(ctx context.Context, email string) (db.Tbluser, error)
	GetAllUsers(ctx context.Context, arg db.GetAllUsersParams) ([]db.Tbluser, error)
	GetUsersWithFilters(ctx context.Context, arg db.GetUsersWithFiltersParams) ([]db.Tbluser, error)
}

// SQLUserRepository es la implementación de UserRepository usando SQLC.
type SQLUserRepository struct {
	*db.Queries
	database *sql.DB
}

// NewSQLUserRepository crea una nueva instancia de SQLUserRepository.
func NewSQLUserRepository(database *sql.DB) UserRepository {
	// La conexión 'database' (*sql.DB) satisface la interfaz DBTX
	// que la función db.New() espera, por lo que se puede pasar directamente.
	return &SQLUserRepository{
		Queries:  db.New(database),
		database: database,
	}
}
