package router

import (
	"challange10/controllers"
	"challange10/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	productRouter := r.Group("/products")
	{
		productRouter.Use(middlewares.Authentication())
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.GET("/", middlewares.ProductAuthorization(), middlewares.AdminAuthorization(), controllers.GetProducts)
		productRouter.GET("/:productId", middlewares.ProductAuthorization(), controllers.GetProduct)
		productRouter.PUT("/:productId", middlewares.ProductAuthorization(), middlewares.AdminAuthorization(), controllers.UpdateProduct)
		productRouter.DELETE("/:productId", middlewares.ProductAuthorization(), middlewares.AdminAuthorization(), controllers.DeleteProduct)
	}

	return r
}
