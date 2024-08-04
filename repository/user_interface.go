package repository

import (
	"github.com/google/uuid"
	"github.com/markmumba/project-tracker/models"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByEmail(email string, user *models.User) error
	GetUser(id uuid.UUID) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetStudentsByLecturer(lecturerID uuid.UUID) ([]models.User, error)
	GetLecturers() ([]models.User, error)
	UpdateUser(id uuid.UUID, user *models.User) error
	UpdateUserProfileImage(id uuid.UUID, profileImage string) error
	DeleteUser(id uuid.UUID) error
}
