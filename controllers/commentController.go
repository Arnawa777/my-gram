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

func commentMap(comment models.Comment) gin.H {
	return gin.H{
		"id":         comment.ID,
		"user_id":    comment.UserID,
		"photo_id":   comment.PhotoID,
		"message":    comment.Message,
		"created_at": comment.CreatedAt,
		"updated_at": comment.UpdatedAt,
	}
}

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

	output := commentMap(comment)
	ctx.JSON(http.StatusOK, output)
}

func CreateComment(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	comment := models.Comment{}

	err := ctx.ShouldBindJSON(&comment)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	comment.UserID = uint(userData["id"].(float64))

	if errs := models.GetValidationErrors(comment); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	err = db.Create(&comment).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	output := commentMap(comment)

	ctx.JSON(http.StatusCreated, output)
}

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

	output := commentMap(comment)
	ctx.JSON(http.StatusOK, output)
}

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

	ctx.Status(http.StatusOK)
}
