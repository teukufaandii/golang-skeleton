package services

import (
	"context"
	"errors"
	"fmt"
	"golang-skeleton/config"
	"golang-skeleton/internal/dto/request"
	"golang-skeleton/internal/dto/response"
	"golang-skeleton/internal/models"
	"golang-skeleton/internal/repository"
	"golang-skeleton/pkg/utils"
	"time"
)

var (
	ErrEmailExists        = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotFound       = errors.New("default role not found")
)

type AuthService interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
}

type authService struct {
	userRepo         repository.UserRepository
	redisRepo        repository.RedisRepository
	jwtSecret        string
	jwtAccessExpiry  time.Duration
	jwtRefreshExpiry time.Duration
}

func NewAuthService(userRepo repository.UserRepository, redisRepo repository.RedisRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo:         userRepo,
		redisRepo:        redisRepo,
		jwtSecret:        cfg.JWTSecret,
		jwtAccessExpiry:  cfg.JWTAccessExpiry,
		jwtRefreshExpiry: cfg.JWTRefreshExpiry,
	}
}

func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error) {
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email availability: %w", err)
	}
	if exists {
		return nil, ErrEmailExists
	}

	role, err := s.userRepo.FindRoleByCode(ctx, "user")
	if err != nil {
		return nil, ErrRoleNotFound
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   role.ID,
		Role:     *role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return s.generateTokens(ctx, user)
}

func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(ctx, user)
}

func (s *authService) generateTokens(ctx context.Context, user *models.User) (*response.AuthResponse, error) {
	accessToken, err := utils.GenerateToken(user.ID, user.Email, user.Role.Code, s.jwtSecret, s.jwtAccessExpiry)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateToken(user.ID, user.Email, user.Role.Code, s.jwtSecret, s.jwtRefreshExpiry)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// create user session token and save it in Redis
	sessionKey := fmt.Sprintf("session:user:%s:%s", user.ID.String(), refreshToken)
	if err := s.redisRepo.Set(ctx, sessionKey, refreshToken, s.jwtRefreshExpiry); err != nil {
		return nil, errors.New("failed to save session to redis")
	}

	return &response.AuthResponse{
		User:         *response.ToUserResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.jwtAccessExpiry).Unix(),
	}, nil
}