package response

import (
	"golang-skeleton/internal/models"
	"time"

	"github.com/google/uuid"
)

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    int64        `json:"expires_at"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role.Code,
		CreatedAt: user.CreatedAt,
	}
}
