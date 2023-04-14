package router

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("users")
	{
		userRouter.POST("/register", controllers.RegisterUser)
		userRouter.POST("/login", controllers.LoginUser)
	}

	photoRouter := r.Group("photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.GET("/", controllers.GetAllPhotos)
		photoRouter.GET("/:ID", controllers.GetPhotoById)
		photoRouter.POST("/", controllers.CreatePhoto)
		photoRouter.PUT("/:ID", middlewares.CheckID("Photo"), controllers.UpdatePhoto)    //need more authorization
		photoRouter.DELETE("/:ID", middlewares.CheckID("Photo"), controllers.DeletePhoto) //need more authorization
		photoRouter.GET("/:ID/comments", controllers.GetAllComment)
		// photoRouter.GET("/:ID/comments", middlewares.ValidateAccess("Photo"), controllers.GetComment)
	}

	commentRouter := r.Group("comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.GET("/:ID", controllers.GetCommentById)
		commentRouter.POST("/:photoID", controllers.CreateComment)
		commentRouter.PUT("/:ID", middlewares.CheckID("Comment"), controllers.UpdateComment)    //need more authorization
		commentRouter.DELETE("/:ID", middlewares.CheckID("Comment"), controllers.DeleteComment) //need more authorization
	}

	sosmedRouter := r.Group("social-media")
	{
		sosmedRouter.Use(middlewares.Authentication())
		sosmedRouter.GET("/", controllers.GetAllSocilaMedia)
		sosmedRouter.GET("/:ID", controllers.GetSocialMediaById)
		sosmedRouter.POST("/", controllers.CreateSocialMedia)
		sosmedRouter.PUT("/:ID", middlewares.CheckID("SocialMedia"), controllers.UpdateSocialMedia)    //need more authorization
		sosmedRouter.DELETE("/:ID", middlewares.CheckID("SocialMedia"), controllers.DeleteSocialMedia) //need more authorization
	}

	return r
}
