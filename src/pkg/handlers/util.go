package handlers

import (
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
)

func isUserSignedIn(ctx *gin.Context) bool {
	user := getAuthenticatedUser(ctx)
	return user != nil && time.Now().Before(user.ExpiresAt)
}

func getAuthenticatedUser(ctx *gin.Context) *data.User {
	session := sessions.Default(ctx)
	user, ok := session.Get(localsUserKey).(data.User)
	if ok {
		return &user
	}
	return nil
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
