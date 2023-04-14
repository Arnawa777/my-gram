package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Create user
// @Description Create new user
// @Tags json
// @Accept json
// @Produce json
// @Param models.User body models.User true "Create User"
// @Success 201 {object} models.User
// @Router /users/register [post]
func RegisterUser(ctx *gin.Context) {
	db := database.GetDB()
	user := models.User{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Check input was correct with validation
	if errs := models.GetValidationErrors(user); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	err = db.Create(&user).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"age":      user.Age,
	})
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticate a user and generate a JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param email body string true "Email address"
// @Param password body string true "Password"
// @Success 200 {object} TokenOutput
// @Failure 400 {object} ErrorOutput
// @Failure 500 {object} ErrorOutput
// @Router /users/login [post]
func LoginUser(ctx *gin.Context) {
	db := database.GetDB()
	user := models.User{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	password := user.Password

	// err to get error from email when it's wrong
	err = db.Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": "invalid email or password",
		})
		return
	}

	// Password Valid return false do this code
	// to give error message and status
	if !helpers.PasswordValid(user.Password, password) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": "invalid email or password",
		})
		return
	}

	// initial token with err
	token, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

type TokenOutput struct {
	Token string `json:"token"`
}

type ErrorOutput struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
