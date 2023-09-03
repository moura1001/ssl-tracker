package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func HandleGetHome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "home/index", util.Map{})
}
