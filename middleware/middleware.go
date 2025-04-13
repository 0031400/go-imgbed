package middleware

import (
	"imgbed/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(c config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token != c.Server.Token {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 1, "data": "not allowed"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
