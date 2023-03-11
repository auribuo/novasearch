package routes

import "github.com/gin-gonic/gin"

func CalculateGalaxyPath(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}

func CalculateViewportPath(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}
