package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
)

func isUserSignedIn(ctx *gin.Context) bool {
	user := getAuthenticatedUser(ctx)
	return user != nil
}

func getAuthenticatedUser(ctx *gin.Context) *data.User {
	value, exist := ctx.Get(localsUserKey)
	if exist {
		if user, ok := value.(*data.User); ok {
			return user
		}
	}
	return nil
}
