package repository

import (
	"github.com/google/uuid"
	"github.com/markmumba/project-tracker/models"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetProject(id uuid.UUID) (*models.Project, error)
	GetProjectsByLecturerId(lecturerId uuid.UUID) ([]models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id uuid.UUID) error
}
