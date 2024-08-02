package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/markmumba/project-tracker/auth"
	"github.com/markmumba/project-tracker/helpers"
	"github.com/markmumba/project-tracker/models"
	"github.com/markmumba/project-tracker/services"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) Login(c echo.Context) error {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	accessToken, refreshToken, err := uc.UserService.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	// Set refresh token as HTTP-only cookie
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Return access token in the response body
	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
	})
}

func (uc *UserController) Refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Missing refresh token")
	}

	claims, err := auth.ValidateRefreshToken(refreshToken.Value)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid refresh token")
	}

	var user models.User
	err = uc.UserService.FindUserByID(claims.UserId, &user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "User not found")
	}

	accessToken, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate access token")
	}

	newRefreshToken, err := auth.GenerateRefreshToken(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not generate refresh token")
	}

	// Update refresh token in the database
	err = uc.UserService.UpdateRefreshToken(refreshToken.Value, newRefreshToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not update refresh token")
	}

	// Set new refresh token as HTTP-only cookie
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Return new access token in the response body
	return c.JSON(http.StatusOK, echo.Map{
		"access_token": accessToken,
	})
}

func (uc *UserController) Logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Missing refresh token")
	}

	err = uc.UserService.DeleteRefreshToken(refreshToken.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Could not delete refresh token")
	}

	// Clear the refresh token cookie
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	return c.JSON(http.StatusOK, "Logged out successfully")
}






func (uc *UserController) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := uc.UserService.CreateUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, models.UserToDTO(&user))
}

func (uc *UserController) GetUser(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	user, err := uc.UserService.GetUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, models.UserToDTO(user))
}

func (uc *UserController) GetAllUsers(c echo.Context) error {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.UserToDTOs(users))
}

func (uc *UserController) GetLecturers(c echo.Context) error {
	lecturers, err := uc.UserService.GetLecturers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, models.UserToDTOs(lecturers))
}

func (uc *UserController) GetStudentsByLecturerId(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	students, err := uc.UserService.GetStudentsByLecturer(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, models.UserToDTOs(students))
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	var updateUser models.User
	if err := c.Bind(&updateUser); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Printf("this is the data from the user %v",updateUser)

	if err := uc.UserService.UpdateUser(userID, &updateUser); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, models.UserToDTO(&updateUser))
}

func (uc *UserController) UpdateUserProfileImage(c echo.Context) error {
	var image struct {
		ProfileImage string `json:"profile_image"`
	}
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Bind(&image); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = uc.UserService.UpdateUserProfileImage(userID, image.ProfileImage)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Profile image updated successfully")
}

func (uc *UserController) DeleteUser(c echo.Context) error {
	userID, err := helpers.ConvertUserID(c, "userId")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	err = uc.UserService.DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, "User deleted successfully")
}
