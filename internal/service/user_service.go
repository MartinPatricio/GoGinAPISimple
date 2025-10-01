package service

import (
	"context"
	"errors"
	"time"

	"github.com/MartinPatricio/GoGinAPISimple/internal/config"
	"github.com/MartinPatricio/GoGinAPISimple/internal/repository"
	"github.com/MartinPatricio/GoGinAPISimple/internal/repository/db"
	"github.com/MartinPatricio/GoGinAPISimple/pkg/hash"
	"github.com/MartinPatricio/GoGinAPISimple/pkg/token"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	req.LastActivitie = pgtype.Date{Time: time.Now(), Valid: true}
	req.DateCreated = pgtype.Date{Time: time.Now(), Valid: true}
	return s.repo.CreateUser(ctx, req)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	jwtToken, err := token.GenerateToken(user.IdUser, s.config.JwtSecretKey, s.config.JwtExpirationHours)
	if err != nil {
		return "", err
	}
	return jwtToken, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, arg db.GetAllUsersParams) ([]db.Tbluser, error) {
	return s.repo.GetAllUsers(ctx, arg)
}

func (s *UserService) GetUsersWithFilters(ctx context.Context, arg db.GetUsersWithFiltersParams) ([]db.Tbluser, error) {
	return s.repo.GetUsersWithFilters(ctx, arg)
}

// internal/service/user_service.go

// ... (resto de tu código de servicio) ...

func (s *UserService) GetUserByID(ctx context.Context, id int32) (db.Tbluser, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}
