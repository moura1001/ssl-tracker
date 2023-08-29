package handlers

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func flashWithData(ctx *gin.Context, value util.Map) {
	session := sessions.Default(ctx)
	session.Set("flash", value)
	if err := session.Save(); err != nil {
		logger.Log("error", "flash error", "err", fmt.Errorf("error in flashMessage saving session: %s", err))
	}
}

func flashes(ctx *gin.Context) util.Map {
	session := sessions.Default(ctx)
	flashes, _ := session.Get("flash").(util.Map)
	if len(flashes) > 0 {
		session.Set("flash", util.Map{})
		if err := session.Save(); err != nil {
			logger.Log("error", "flash error", "err", fmt.Errorf("error in flashes saving session: %s", err))
		}
	}
	return flashes
}
