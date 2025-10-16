package main

import (
	"fmt"

	"StreamMovieServer/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// This is main function
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello!!, StreamMovieServer")
	})

	routes.SetupUnProtectedRoutes(router)
	routes.SetupProtectedRoutes(router)

	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)

	}
}
