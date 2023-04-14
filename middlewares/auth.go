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

func getModel(table string) (interface{}, error) {
	switch table {
	case "Photo":
		return models.Photo{}, nil
	case "SocialMedia":
		return models.SocialMedia{}, nil
	case "Comment":
		return models.Comment{}, nil
	default:
		return nil, fmt.Errorf("Invalid table %s", table)
	}
}

func Authorization(table string) gin.HandlerFunc {
	tables := map[string]interface{}{
		"Photo":       models.Photo{},
		"SocialMedia": models.SocialMedia{},
		"Comment":     models.Comment{},
	}

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
		model, ok := tables[table]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "Unauthorized",
				"message": fmt.Sprintf("Invalid table %s", table),
			})
			return
		}
		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))

		res := db.Model(model).Select("user_id").First(model, uint(ID))

		var domainUserID uint
		switch v := model.(type) {
		case *models.Photo:
			domainUserID = v.UserID
		case *models.SocialMedia:
			domainUserID = v.UserID
		case *models.Comment:
			domainUserID = v.UserID
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
