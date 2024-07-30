package repository

import (

	"github.com/markmumba/project-tracker/database"
	"github.com/markmumba/project-tracker/models"
)



type RefreshTokenRepositoryImpl struct{}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{}
}

func (r *RefreshTokenRepositoryImpl) Create(token *models.RefreshToken) error {
	return database.DB.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) SaveRefreshToken(token *models.RefreshToken) error {
    return database.DB.Save(token).Error
}

func (r *RefreshTokenRepositoryImpl) FindRefreshToken(token string, refreshToken *models.RefreshToken) error {
    return database.DB.Where("token = ?", token).First(refreshToken).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteRefreshToken(token string) error {
    return database.DB.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}