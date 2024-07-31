package services

import (
	"errors"
	"time"

	"github.com/markmumba/project-tracker/auth"
	"github.com/markmumba/project-tracker/models"
	"github.com/markmumba/project-tracker/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
	RefreshRepository repository.RefreshTokenRepository
}

func NewUserService(userRepo repository.UserRepository,refreshRepo repository.RefreshTokenRepository) *UserService {
	return &UserService{
		UserRepository: userRepo,
		RefreshRepository: refreshRepo,
	}
}

func (u *UserService) CreateUser(user *models.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return u.UserRepository.CreateUser(user)
}

func (u *UserService) LoginUser(email, password string) (string, string, error) {
    var user models.User
    err := u.UserRepository.FindByEmail(email, &user)
    if err != nil {
        return "", "", err
    }

    if !auth.CheckPasswordHash(password, user.Password) {
        return "", "", errors.New("invalid credentials")
    }

    accessToken, err := auth.GenerateAccessToken(&user)
    if err != nil {
        return "", "", err
    }

    refreshToken, err := auth.GenerateRefreshToken(&user)
    if err != nil {
        return "", "", err
    }

    // Store refresh token in the database
    err = u.RefreshRepository.Create(&models.RefreshToken{
        Token:     refreshToken,
        UserID:    user.ID,
        ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
    })
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

func (u *UserService) FindRefreshToken(token string, refreshToken *models.RefreshToken) error {
    return u.RefreshRepository.FindRefreshToken(token, refreshToken)
}

func (u *UserService) FindUserByID(userID uint, user *models.User) error {
    return u.FindUserByID(userID, user)
}

func (u *UserService) UpdateRefreshToken(oldToken, newToken string) error {
	return u.RefreshRepository.UpdateToken(oldToken, newToken)
}

func (u *UserService) SaveRefreshToken(token *models.RefreshToken) error {
    return u.RefreshRepository.SaveRefreshToken(token)
}

func (u *UserService) DeleteRefreshToken(token string) error {
    return u.RefreshRepository.DeleteRefreshToken(token)
}

func (u *UserService) GetUser(id uint) (*models.User, error) {
	return u.UserRepository.GetUser(id)
}

func (u *UserService) GetAllUsers() ([]models.User, error) {
	return u.UserRepository.GetAllUsers()
}

func (u *UserService) GetStudentsByLecturer(lecturerID uint) ([]models.User, error) {
	return u.UserRepository.GetStudentsByLecturer(lecturerID)
}

func (u *UserService) GetLecturers() ([]models.User, error) {
	return u.UserRepository.GetLecturers()
}

func (u *UserService) UpdateUser(id uint, user *models.User) error {

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return u.UserRepository.UpdateUser(id, user)
}

func (u *UserService) UpdateUserProfileImage(id uint, profileImage string) error {
	return u.UserRepository.UpdateUserProfileImage(id, profileImage)
}

func (u *UserService) DeleteUser(id uint) error {
	return u.UserRepository.DeleteUser(id)
}