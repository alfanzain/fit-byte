package routes

import (
	"database/sql"
	"fit-byte/config"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, db *sql.DB) *gin.Engine {
	router := gin.Default()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("isImage", utils.IsImageURI)
	// }

	// v1Group := router.Group("/api/v1")
	// {
	// v1Group.POST("/login", AuthHandler.Login)
	// v1Group.POST("/post", AuthHandler.Register)

	// v1Group.GET("/user", UserHandler.GetUsers)
	// v1Group.PATCH("/user", UserHandler.UpdateUser)

	// v1Group.POST("/file", FileHandler.UploadFile)

	// v1Group.POST("/activity", ActivityHandler.CreateActivity)
	// v1Group.GET("/activity", ActivityHandler.GetActivities)
	// v1Group.PATCH("/activity/:activityId", ActivityHandler.UpdateActivity)
	// v1Group.DELETE("/activity/:activityId", ActivityHandler.DeleteActivity)
	// }

	return router
}
