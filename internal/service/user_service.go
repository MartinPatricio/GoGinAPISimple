package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MartinPatricio/GoGinAPISimple/internal/repository"
	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
	"github.com/MartinPatricio/GoGinAPISimple/pkg/hash"
	"github.com/MartinPatricio/GoGinAPISimple/pkg/token"

	"github.com/MartinPatricio/GoGinAPISimple/internal/config"
)

type UserService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:   repo,
		config: cfg,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req db.CreateUserParams) (db.Tbluser, error) {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return db.Tbluser{}, err
	}
	req.Password = hashedPassword
	req.Lastactivitie = sql.NullTime{Time: time.Now(), Valid: true}

	return s.repo.CreateUser(ctx, req)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	jwtToken, err := token.GenerateToken(user.Iduser, s.config.JwtSecretKey, s.config.JwtExpirationHours)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// Añade aquí los demás métodos del servicio que llamarán al repositorio
// (DeleteUser, GetUserByID, GetAllUsers, etc.)
