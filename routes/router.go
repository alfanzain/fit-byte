package routes

import (
	"database/sql"
	"fit-byte/config"
	v1Handlers "fit-byte/handlers/v1"
	"fit-byte/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()
	jwtMiddleware := middleware.JWTAuth()

	v1Group := router.Group("/v1")

	authHandler := v1Handlers.NewAuthHandler(db)
	activityHandler := v1Handlers.NewActivityHandler(db)

	v1Group.POST("/login", authHandler.Login)
	v1Group.POST("/register", authHandler.Register)

	// userRouter := v1Group.Group("user")
	// userRouter.Use(jwtMiddleware)
	// userRouter.GET("/", Handler.GetUsers)
	// userRouter.PATCH("/:userId", UserHandler.UpdateUser)

	// fileRouter := v1Group.Group("file")
	// fileRouter.Use(jwtMiddleware)
	// fileRouter.POST("/", FileHandler.UploadFile)

	activityRouter := v1Group.Group("activity")
	activityRouter.Use(jwtMiddleware)
	activityRouter.POST("/", activityHandler.CreateActivity)
	activityRouter.GET("/", activityHandler.GetActivities)
	activityRouter.PATCH("/:activityId", activityHandler.UpdateActivity)
	activityRouter.DELETE("/:activityId", activityHandler.DeleteActivity)

	testRouter := v1Group.Group("middleware-test")
	testRouter.Use(jwtMiddleware)
	testRouter.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "middleware works",
		})
	})

	return router
}
