package routes

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markmumba/project-tracker/controllers"
	"github.com/markmumba/project-tracker/custommiddleware"
	"github.com/markmumba/project-tracker/services"
)

var vercelFrontend = os.Getenv("VERCEL_FRONTEND")
var railwayFrontend = os.Getenv("RAILWAY_FRONTEND")

func SetupRouter(
	userService *services.UserService,
	projectService *services.ProjectService,
	submissionService *services.SubmissionService,
	feedbackService *services.FeedbackService,
	communicationService *services.CommunicationService,
) *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", vercelFrontend, railwayFrontend},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
		MaxAge:           300, // Optional: cache preflight requests
	}))

	userController := controllers.NewUserController(userService)
	projectController := controllers.NewProjectController(projectService)
	submissionController := controllers.NewSubmissionController(submissionService)
	feedbackController := controllers.NewFeedbackController(feedbackService)
	communicationController := controllers.NewCommunicationContoller(communicationService)
	websocketController := controllers.NewWebsocketController(*communicationService, *projectService)

	e.POST("/register", userController.CreateUser)
	e.POST("/login", userController.Login)
	e.GET("/logout", userController.Logout)
	e.GET("/ws", websocketController.HandleWebSocket)
	e.GET("/refresh-token", userController.Refresh)

	r := e.Group("")
	r.Use(custommiddleware.Authentication)

	userGroup := r.Group("/users")
	{
		userGroup.GET("", userController.GetUser)
		userGroup.GET("/all", userController.GetAllUsers)
		userGroup.GET("/students", userController.GetStudentsByLecturerId)
		userGroup.GET("/lecturers", userController.GetLecturers)
		userGroup.PUT("", userController.UpdateUser)
		userGroup.POST("/profile", userController.UpdateUserProfileImage)
		userGroup.DELETE("", userController.DeleteUser)
	}

	projectGroup := r.Group("/projects")
	{
		projectGroup.POST("", projectController.CreateProject)
		projectGroup.GET("", projectController.GetProject)
		projectGroup.PUT("", projectController.UpdateProject)
		projectGroup.DELETE("", projectController.DeleteProject)
	}

	submissionGroup := r.Group("/submissions")
	{
		submissionGroup.POST("", submissionController.CreateSubmission)
		submissionGroup.GET("/:id", submissionController.GetSubmission)
		submissionGroup.GET("/student", submissionController.GetAllSubmissionByStudentId)
		submissionGroup.GET("/lecturer", submissionController.GetSubmissionsByLecturer)
		submissionGroup.PUT("/:id", submissionController.UpdateSubmission)
		submissionGroup.GET("", submissionController.GetAllSubmissions)
		submissionGroup.DELETE("/:id", submissionController.DeleteSubmission)
	}

	feedbackGroup := r.Group("/feedbacks")
	{
		feedbackGroup.POST("", feedbackController.CreateFeedback)
		feedbackGroup.GET("/student", feedbackController.GetFeedbackByStudent)
		feedbackGroup.GET("/lecturer", feedbackController.GetFeedbackByLecturer)
		feedbackGroup.GET("", feedbackController.GetAllFeedback)
		feedbackGroup.PUT("/:id", feedbackController.UpdateFeedback)
		feedbackGroup.DELETE("/:id", feedbackController.DeleteFeedback)
		feedbackGroup.GET("/submission/:id", feedbackController.GetFeedbackForSubmission)
	}

	communicationGroup := r.Group("/communications")
	{
		communicationGroup.POST("", communicationController.SaveMessage)
		communicationGroup.GET("", communicationController.GetMessagesBetweenUsers)
	}

	return e
}
