package controllers

import (
	"final-project/database"
	"final-project/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// GetAllComment godoc
// @Summary Get all comments
// @Description Get details about all comments
// @Tags json
// @Accept json
// @Produce json
// @Success 200 {object} models.Comment
// @Router /comments [get]
func GetAllComment(ctx *gin.Context) {
	db := database.GetDB()
	comments := []models.Comment{}
	ID, err := strconv.Atoi(ctx.Param("ID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": err.Error(),
		})
		return
	}

	result := db.Where("photo_id = ?", ID).Find(&comments)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// GetCommentById godoc
// @Summary Get comment by comment id
// @Description Get details of specific comment
// @Tags json
// @Accept json
// @Produce json
// @Param Id path uint true "ID of the comment"
// @Success 200 {object} models.Comment
// @Router /comments/{Id} [get]
func GetCommentById(ctx *gin.Context) {
	db := database.GetDB()
	comment := models.Comment{}
	commentID, _ := strconv.Atoi(ctx.Param("ID"))

	err := db.First(&comment, commentID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// CreateComment godoc
// @Summary Create comment
// @Description Create new comment
// @Tags json
// @Accept json
// @Produce json
// @Param models.Comment body models.Comment true "Create comment"
// @Success 201 {object} models.Comment
// @Router /comments [post]
func CreateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	comment := models.Comment{}
	photoID, err := strconv.Atoi(ctx.Param("photoID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": err.Error(),
		})
		return
	}

	err = ctx.ShouldBindJSON(&comment)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	comment.UserID = uint(userData["id"].(float64))
	comment.PhotoID = uint(photoID)

	if errs := models.GetValidationErrors(comment); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	err = db.Create(&comment).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

// UpdateComment godoc
// @Summary Update comment
// @Description Update comment data
// @Tags json
// @Accept json
// @Produce json
// @Param Id path int true "ID of the comment"
// @Success 200 {object} models.Comment
// @Router /comments/{Id} [patch]
func UpdateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	comment := models.Comment{}
	commentID, _ := strconv.Atoi(ctx.Param("ID"))

	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	newMessage := comment.Message

	// Fetch the comment from the database
	var dbComment models.Comment
	err = db.First(&dbComment, commentID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	comment = dbComment
	comment.UserID = uint(userData["id"].(float64))
	comment.Message = newMessage

	// Update Comment
	result := db.Model(&comment).Where("id=?", commentID).Updates(models.Comment{Message: newMessage})
	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	// Validate before Update
	if errs := models.GetValidationErrors(comment); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Comment with ID %d not found", commentID),
		})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

// DeleteComment godoc
// @Summary Delete comment
// @Description Delete comment data
// @Tags json
// @Accept json
// @Produce json
// @Param Id path int true "ID of the comment"
// @Success 200
// @Router /comments/{Id} [delete]
func DeleteComment(ctx *gin.Context) {
	db := database.GetDB()
	comment := models.Comment{}
	commentID, err := strconv.Atoi(ctx.Param("ID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": err.Error(),
		})
		return
	}

	// Delete the comment from the database
	if err := db.Delete(&comment, uint(commentID)).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success to delete comment",
	})
}
