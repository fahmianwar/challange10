package controllers

import (
	"net/http"
	"strconv"

	"challange10/database"
	"challange10/helpers"
	"challange10/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GetProducts(ctx *gin.Context) {
	db := database.GetDB()
	products := []models.Product{}
	db.Find(&products)

	for i := 0; i < len(products); i++ {
		if err := db.Debug().Where("id = ?", products[i].UserID).First(&products[i].User).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": products})
}

func GetProduct(ctx *gin.Context) {
	product := models.Product{}

	db := database.GetDB()
	if err := db.Where("id = ?", ctx.Param("productId")).First(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Debug().Where("id = ?", product.UserID).First(&product.User).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(ctx *gin.Context) {
	db := database.GetDB()
	product := models.Product{}
	if err := db.Where("id = ?", ctx.Param("productId")).First(&product).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Data not found!!"})
		return
	}

	db.Delete(&product)

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted Product Success"})
}

func CreateProduct(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContenType(ctx)

	product := models.Product{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&product)
	} else {
		ctx.ShouldBind(&product)
	}

	product.UserID = userID

	if err := db.Debug().Where("id = ?", userID).First(&product.User).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found!"})
		return
	}
	err := db.Debug().Create(&product).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Created Product Success", "data": product})
}

func UpdateProduct(ctx *gin.Context) {
	db := database.GetDB()
	userData := ctx.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContenType(ctx)
	product := models.Product{}

	productId, _ := strconv.Atoi(ctx.Param("productId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&product)
	} else {
		ctx.ShouldBind(&product)
	}

	product.UserID = userID
	product.ID = uint(productId)

	err := db.Model(&product).Where("id = ?", productId).Updates(models.Product{Title: product.Title, Description: product.Description}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}
