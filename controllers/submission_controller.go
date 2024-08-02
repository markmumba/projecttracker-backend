package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/markmumba/project-tracker/helpers"
	"github.com/markmumba/project-tracker/models"
	"github.com/markmumba/project-tracker/services"
)

// TODO : streamline type conversion from frontend to backend and vice versa
type SubmissionController struct {
	submissionService *services.SubmissionService
}

func NewSubmissionController(submissionService *services.SubmissionService) *SubmissionController {
	return &SubmissionController{
		submissionService: submissionService,
	}
}

func (sc *SubmissionController) CreateSubmission(c echo.Context) error {
	var submission models.Submission
	if err := c.Bind(&submission); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := sc.submissionService.CreateSubmission(&submission); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, models.SubmissionToDTO(&submission))
}

func (sc *SubmissionController) GetSubmission(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	submission, err := sc.submissionService.GetSubmission(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, models.SubmissionToDTO(submission))
}

func (sc *SubmissionController) GetAllSubmissionByStudentId(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	submissions, err := sc.submissionService.GetAllSubmissionByStudentId(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.SubmissionToDTOs(submissions))
}

func (sc *SubmissionController) GetSubmissionsByLecturer(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	submissions, err := sc.submissionService.GetSubmissionsByLecturer(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.SubmissionToDTOs(submissions))
}

func (sc *SubmissionController) UpdateSubmission(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	var submission models.Submission
	if err := c.Bind(&submission); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := sc.submissionService.UpdateSubmission(&submission, id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.SubmissionToDTO(&submission))
}

func (sc *SubmissionController) DeleteSubmission(c echo.Context) error {

	idParam := c.Param("id")
	submissionId, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = sc.submissionService.DeleteSubmission(submissionId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, "Submission deleted successfully")
}

func (sc *SubmissionController) GetAllSubmissions(c echo.Context) error {
	submissions, err := sc.submissionService.GetAllSubmissions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, submissions)
}
