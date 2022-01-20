package api

import (
	"net/http"

	sess "github.com/E_learning/sessions"
	"github.com/gin-gonic/gin"
)

func sessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken := sess.SessionStart().Get("username", ctx)
		if sessionToken == nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Not logged in"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
