package helpers

import "github.com/gin-gonic/gin"

func ResponseErr(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"error": message})
}
