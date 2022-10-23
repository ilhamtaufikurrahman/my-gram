package router

import (
	"my-gram/controllers"
	"my-gram/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.PUT("/:userId", middlewares.Authentication(), controllers.UserUpdate)
	}

	socialMediaRouter := r.Group("/socialmedias")
	{
		socialMediaRouter.Use(middlewares.Authentication())
		socialMediaRouter.PUT("/:socialMediaId", controllers.UpdateSocialMedia)
		socialMediaRouter.DELETE("/:socialMediaId", controllers.DeleteSocialMedia)
		socialMediaRouter.POST("/", controllers.CreateSocialMedia)
		socialMediaRouter.GET("/", controllers.GetSocialMedias)
	}

	photosRouter := r.Group("/photos")
	{
		photosRouter.Use(middlewares.Authentication())
		photosRouter.POST("/", controllers.CreatePhoto)
		photosRouter.GET("/", controllers.GetPhotos)
		photosRouter.PUT("/:photoId", controllers.UpdatePhoto)
	}

	commentsRouter := r.Group("/comments")
	{
		commentsRouter.Use(middlewares.Authentication())
		commentsRouter.POST("/", controllers.CreateComment)
		commentsRouter.GET("/", controllers.GetComments)
		commentsRouter.PUT("/:commentId", controllers.UpdateComment)
	}

	return r
}
