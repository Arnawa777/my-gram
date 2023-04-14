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

// GetAllPhotos godoc
// @Summary Get all photos
// @Description Get details about all photos
// @Tags json
// @Accept json
// @Produce json
// @Success 200 {object} models.Photo
// @Router /photos [get]
func GetAllPhotos(ctx *gin.Context) {
	db := database.GetDB()
	photos := []models.Photo{}

	err := db.Find(&photos).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, photos)
}

// GetPhotoById godoc
// @Summary Get photo by photo id
// @Description Get details of specific photo
// @Tags json
// @Accept json
// @Produce json
// @Param Id path uint true "ID of the photo"
// @Success 200 {object} models.Photo
// @Router /photos/{Id} [get]
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

	ctx.JSON(http.StatusOK, photo)
}

// CreatePhoto godoc
// @Summary Create photo
// @Description Create new photo
// @Tags json
// @Accept json
// @Produce json
// @Param models.Photo body models.Photo true "Create photo"
// @Success 201 {object} models.Photo
// @Router /photos [post]
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

	ctx.JSON(http.StatusCreated, photo)
}

// UpdatePhoto godoc
// @Summary Update photo
// @Description Update photo data
// @Tags json
// @Accept json
// @Produce json
// @Param Id path int true "ID of the photo"
// @Success 200 {object} models.Photo
// @Router /photos/{Id} [patch]
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

	ctx.JSON(http.StatusOK, photo)
}

// DeletePhoto godoc
// @Summary Delete photo
// @Description Delete photo data
// @Tags json
// @Accept json
// @Produce json
// @Param Id path int true "ID of the photo"
// @Success 200
// @Router /photos/{Id} [delete]
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

	// ctx.Status(http.StatusOK)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "success to delete photo",
	})
}
