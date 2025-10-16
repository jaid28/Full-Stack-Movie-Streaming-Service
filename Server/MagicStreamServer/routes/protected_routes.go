package routes

import (
	controller "StreamMovieServer/controllers"
	"StreamMovieServer/middleware"

	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())

	router.GET("/movies/:imdb_id", controller.GetMovie())
	router.POST("/addmovie", controller.AddMovie())
}
