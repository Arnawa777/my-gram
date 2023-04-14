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

func socialMediaMap(socialMedia models.SocialMedia) gin.H {
	return gin.H{
		"id":               socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaURL,
		"user_id":          socialMedia.UserID,
		"created_at":       socialMedia.CreatedAt,
		"updated_at":       socialMedia.UpdatedAt,
	}
}

func GetAllSocilaMedia(ctx *gin.Context) {
	db := database.GetDB()
	socialMedia := []models.SocialMedia{}
	userData := ctx.MustGet("userData").(jwt.MapClaims)

	result := db.Where("user_id = ?", uint(userData["id"].(float64))).Find(&socialMedia)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": result.Error,
		})
		return
	}

	var output []gin.H
	for _, sosmed := range socialMedia {
		output = append(output, socialMediaMap(sosmed))
	}

	ctx.JSON(http.StatusOK, output)
}

func GetSocialMediaById(ctx *gin.Context) {
	db := database.GetDB()
	sosmed := models.SocialMedia{}
	sosmedID, _ := strconv.Atoi(ctx.Param("ID"))

	err := db.First(&sosmed, sosmedID).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusText(http.StatusInternalServerError),
			"message": err.Error(),
		})
		return
	}

	output := socialMediaMap(sosmed)
	ctx.JSON(http.StatusOK, output)
}

func CreateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	sosmed := models.SocialMedia{}

	err := ctx.ShouldBindJSON(&sosmed)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sosmed.UserID = uint(userData["id"].(float64))

	if errs := models.GetValidationErrors(sosmed); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	err = db.Create(&sosmed).Error
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	output := socialMediaMap(sosmed)

	ctx.JSON(http.StatusCreated, output)
}

func UpdateSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	sosmed := models.SocialMedia{}
	sosmedID, _ := strconv.Atoi(ctx.Param("ID"))
	err := ctx.ShouldBindJSON(&sosmed)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	sosmed.UserID = uint(userData["id"].(float64))
	//just to make output id not 0
	sosmed.ID = uint(sosmedID)

	// Update sosmed
	result := db.Model(&sosmed).Where("id=?", sosmedID).Updates(models.SocialMedia{Name: sosmed.Name, SocialMediaURL: sosmed.SocialMediaURL})
	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, result.Error)
		return
	}

	// Validate before Update
	if errs := models.GetValidationErrors(sosmed); len(errs) > 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Social Media with ID %d not found", sosmedID),
		})
		return
	}

	output := socialMediaMap(sosmed)
	ctx.JSON(http.StatusOK, output)
}

func DeleteSocialMedia(ctx *gin.Context) {
	db := database.GetDB()
	sosmed := models.SocialMedia{}
	sosmedID, err := strconv.Atoi(ctx.Param("ID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusText(http.StatusBadRequest),
			"message": err.Error(),
		})
		return
	}

	// Delete the sosmed from the database
	if err := db.Delete(&sosmed, uint(sosmedID)).Error; err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusOK)
}
