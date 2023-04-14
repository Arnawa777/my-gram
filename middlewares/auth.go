package middlewares

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := helpers.VerifyToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"error":   err.Error(),
			})
			return
		}

		ctx.Set("userData", claims)

		ctx.Next()
	}
}

func CheckID(table string) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		ID, err := strconv.Atoi(c.Param("ID"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Unauthorized",
				"message": "Invalid ID data type",
			})
			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		var res *gorm.DB
		var domainUserID uint

		switch table {
		case "Photo":
			domain := models.Photo{}
			res = db.Select("user_id").First(&domain, uint(ID))
			domainUserID = domain.UserID
		case "SocialMedia":
			domain := models.SocialMedia{}
			res = db.Select("user_id").First(&domain, uint(ID))
			domainUserID = domain.UserID
		case "Comment":
			domain := models.Comment{}
			res = db.Select("user_id").First(&domain, uint(ID))
			domainUserID = domain.UserID
		}

		if res.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": fmt.Sprintf("%s with ID %d not found", table, ID),
			})
			return
		}

		if domainUserID != userID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "Forbidden",
				"message": fmt.Sprintf("You are not allowed to access this %s", table),
			})
			return
		}

		c.Next()
	}
}
