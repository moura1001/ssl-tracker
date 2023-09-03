package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const localsUserKey = "user"

func WithFlash(ctx *gin.Context) {
	values := flashes(ctx)
	ctx.Set("flash", values)
	ctx.Next()
}

func WithViewHelpers(ctx *gin.Context) {
	ctx.Set("activeFor", func(s string) (res string) {
		if ctx.Request.URL.Path == s {
			return "active"
		}
		return ""
	})
	ctx.Next()
}

func WithMustBeAuthenticated(ctx *gin.Context) {
	if !isUserSignedIn(ctx) {
		ctx.Abort()
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	ctx.Next()
}
