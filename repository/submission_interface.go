package repository

import (
	"github.com/google/uuid"
	"github.com/markmumba/project-tracker/models"
)

type SubmissionRepository interface {
	CreateSubmission(submission *models.Submission) error
	GetSubmission(id uuid.UUID) (*models.Submission, error)
	GetAllSubmissionByStudentId(studentId uuid.UUID) ([]models.Submission, error)
	GetSubmissionsByLecturer(lecturerID uuid.UUID) ([]models.Submission, error)
	UpdateSubmission(submission *models.Submission, id uuid.UUID) error
	DeleteSubmission(id uuid.UUID) error
	GetAllSubmissions() ([]models.Submission, error)
}
