package features

import "github.com/gin-gonic/gin"

func AllGalaxies(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}

func Galaxy(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}

func FilterGalaxies(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "pong",
	})
}
