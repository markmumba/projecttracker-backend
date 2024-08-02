package repository

import (
	"github.com/google/uuid"
	"github.com/markmumba/project-tracker/models"
)

type FeedbackRepository interface {
	CreateFeedback(feedback *models.Feedback) error
	GetFeedback(id uuid.UUID) (*models.Feedback, error)
	GetFeedbackByStudent(studentID uuid.UUID) (*[]models.Feedback, error)
	GetFeedbackByLecturer(lecturerID uuid.UUID) (*[]models.Feedback, error)
	GetAllFeedback() ([]models.Feedback, error)
	GetFeedbackBySubmissionId(submissionId uuid.UUID) ([]models.Feedback, error)
	UpdateFeedback(feedback *models.Feedback) error
	GetFeedbackForSubmission(id uuid.UUID) (*models.Feedback, error)
	DeleteFeedback(id uuid.UUID) error
}
