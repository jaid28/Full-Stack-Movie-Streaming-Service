package routes

import (
	controller "StreamMovieServer/controllers"

	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {

	router.GET("/movies", controller.GetMovies())
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
}
