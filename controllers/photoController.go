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

func photoToMap(photo models.Photo) gin.H {
	return gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoURL,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
		"updated_at": photo.UpdatedAt,
	}
}

func GetAllPhotos(ctx *gin.Context) {
	db := database.GetDB()
	photos := []models.Photo{}

	err := db.Find(&photos).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var output []gin.H
	for _, photo := range photos {
		output = append(output, photoToMap(photo))
	}

	ctx.JSON(http.StatusOK, output)
}

func GetPhotoById(ctx *gin.Context) {
	db := database.GetDB()
	photo := models.Photo{}
	photoID, _ := strconv.Atoi(ctx.Param("ID"))

	err := db.First(&photo, photoID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	output := photoToMap(photo)
	ctx.JSON(http.StatusOK, output)
}

func CreatePhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}

	err := ctx.ShouldBindJSON(&photo)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	photo.UserID = uint(userData["id"].(float64))

	if errs := models.GetValidationErrors(photo); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	err = db.Create(&photo).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	output := photoToMap(photo)

	ctx.JSON(http.StatusCreated, output)
}

func UpdatePhoto(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	photo := models.Photo{}
	photoID, _ := strconv.Atoi(ctx.Param("ID"))
	err := ctx.ShouldBindJSON(&photo)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	photo.UserID = uint(userData["id"].(float64))
	//just to make output id not 0
	photo.ID = uint(photoID)

	// Update photo
	result := db.Model(&photo).Where("id=?", photoID).Updates(models.Photo{Title: photo.Title, Caption: photo.Caption, PhotoURL: photo.PhotoURL})
	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	// Validate before Update
	if errs := models.GetValidationErrors(photo); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Photo with ID %d not found", photoID),
		})
		return
	}

	output := photoToMap(photo)
	ctx.JSON(http.StatusOK, output)
}

func DeletePhoto(ctx *gin.Context) {
	db := database.GetDB()
	photo := models.Photo{}
	photoID, err := strconv.Atoi(ctx.Param("ID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": err.Error(),
		})
		return
	}

	// Delete the photo from the database
	if err := db.Delete(&photo, uint(photoID)).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetAllComment(ctx *gin.Context) {
	db := database.GetDB()
	photo := models.Photo{}
	comments := []models.Comment{}
	photoID, _ := strconv.Atoi(ctx.Param("ID"))

	err := db.Preload("Comments").First(&photo, photoID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	comments = photo.Comments // assign preloaded comments to the variable

	output := photoToMap(photo)
	if len(comments) == 0 {
		output["comments"] = []gin.H{}
	} else {
		commentMaps := make([]gin.H, len(comments))
		for i, comment := range comments {
			commentMaps[i] = gin.H{
				"id":         comment.ID,
				"user_id":    comment.UserID,
				"message":    comment.Message,
				"created_at": comment.CreatedAt,
				"updated_at": comment.UpdatedAt,
			}
		}
		output["comments"] = commentMaps
	}

	ctx.JSON(http.StatusOK, output)
}
