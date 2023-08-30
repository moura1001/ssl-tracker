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
	return &data.User{Id: "123", Email: "email@email.com"}
	//return nil
}

func processResults(resultsChan chan data.DomainTracking) []data.DomainTracking {
	var (
		trackings = make([]data.DomainTracking, len(resultsChan))
		i         int
	)
	for result := range resultsChan {
		trackings[i] = result
		i++
	}
	return trackings
}
