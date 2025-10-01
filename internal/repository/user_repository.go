/*package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Tbluser, error)
	DeleteUser(ctx context.Context, idUser int32) error
	GetUserByID(ctx context.Context, idUser int32) (db.Tbluser, error)
	GetUserByEmail(ctx context.Context, email string) (db.Tbluser, error)
	GetAllUsers(ctx context.Context, arg db.GetAllUsersParams) ([]db.Tbluser, error)
	GetUsersWithFilters(ctx context.Context, arg db.GetUsersWithFiltersParams) ([]db.Tbluser, error)
}

type SQLUserRepository struct {
	*db.Queries
	database *pgxpool.Pool
}

func NewSQLUserRepository(database *pgxpool.Pool) UserRepository {
	return &SQLUserRepository{
		Queries:  db.New(database),
		database: database,
	}
}
*/
// internal/repository/user_repository.go
package repository

import (
	"context"
	"log" // <-- 1. Añade el import de log

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Tbluser, error)
	DeleteUser(ctx context.Context, idUser int32) error
	GetUserByID(ctx context.Context, idUser int32) (db.Tbluser, error)
	GetUserByEmail(ctx context.Context, email string) (db.Tbluser, error)
	GetAllUsers(ctx context.Context, arg db.GetAllUsersParams) ([]db.Tbluser, error)
	GetUsersWithFilters(ctx context.Context, arg db.GetUsersWithFiltersParams) ([]db.Tbluser, error)
}

type SQLUserRepository struct {
	*db.Queries
	database *pgxpool.Pool
}

func NewSQLUserRepository(database *pgxpool.Pool) UserRepository {
	return &SQLUserRepository{
		Queries:  db.New(database),
		database: database,
	}
}

// ✅ 2. Añadimos una nueva interfaz para poder registrar las llamadas
// Esto nos permite "envolver" la implementación de SQLC con nuestros logs.
type loggedUserRepository struct {
	next UserRepository
}

func NewLoggedUserRepository(next UserRepository) UserRepository {
	return &loggedUserRepository{
		next: next,
	}
}

// ✅ 3. Implementamos el método CreateUser con logging
func (r *loggedUserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.Tbluser, error) {
	log.Printf("--- Repositorio: Intentando crear usuario con Email: %s ---", arg.Email)

	// Llamamos a la función original de SQLC
	user, err := r.next.CreateUser(ctx, arg)

	if err != nil {
		log.Printf("--- Repositorio: ERROR al crear usuario: %v ---", err)
	} else {
		log.Printf("--- Repositorio: ÉXITO. Usuario creado con ID: %d ---", user.IdUser)
	}

	return user, err
}

// Simplemente pasamos las otras llamadas al siguiente repositorio sin log
func (r *loggedUserRepository) DeleteUser(ctx context.Context, idUser int32) error {
	return r.next.DeleteUser(ctx, idUser)
}
func (r *loggedUserRepository) GetUserByID(ctx context.Context, idUser int32) (db.Tbluser, error) {
	return r.next.GetUserByID(ctx, idUser)
}
func (r *loggedUserRepository) GetUserByEmail(ctx context.Context, email string) (db.Tbluser, error) {
	return r.next.GetUserByEmail(ctx, email)
}
func (r *loggedUserRepository) GetAllUsers(ctx context.Context, arg db.GetAllUsersParams) ([]db.Tbluser, error) {
	return r.next.GetAllUsers(ctx, arg)
}
func (r *loggedUserRepository) GetUsersWithFilters(ctx context.Context, arg db.GetUsersWithFiltersParams) ([]db.Tbluser, error) {
	return r.next.GetUsersWithFilters(ctx, arg)
}
