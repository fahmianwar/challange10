package middlewares

import (
	"net/http"
	"strconv"
	"strings"

	"challange10/database"
	"challange10/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ProductAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		db := database.GetDB()
		productId, err := strconv.Atoi(ctx.Param("productId"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "invalid parameter",
			})
			return
		}
		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		product := models.Product{}

		err = db.Select("user_id").First(&product, uint(productId)).Error

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "Data Not Found",
				"message": "data doesn't exist",
			})
			return
		}

		if product.UserID != userID {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorizad",
				"message": "you are not allowed to access this data",
			})
			return
		}
		ctx.Next()
	}
}

func AdminAuthorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userData := ctx.MustGet("userData").(jwt.MapClaims)
		userLevel := userData["level"].(string)

		if !strings.EqualFold(strings.ToLower(userLevel), strings.ToLower("Admin")) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorizad",
				"message": "you are not allowed to access this data",
			})
			return
		}
		ctx.Next()
	}
}
