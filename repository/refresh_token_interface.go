package repository

import "github.com/markmumba/project-tracker/models"

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	SaveRefreshToken(token *models.RefreshToken) error
	FindRefreshToken(token string, refreshToken *models.RefreshToken) error 
	UpdateToken(oldToken, newToken string) error 
	DeleteRefreshToken(token string) error
}
