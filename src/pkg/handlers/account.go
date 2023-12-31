package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

const maxNotifyUpfront = 356 / 2

type UpdateAccountParams struct {
	NotifyUpfront      int
	DefaultNotifyEmail string
}

func HandleAccountUpdate(ctx *gin.Context) {
	var params UpdateAccountParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	errors := util.Map{}
	if params.NotifyUpfront <= 0 || params.NotifyUpfront > maxNotifyUpfront {
		errors["notifyUpfront"] = fmt.Sprintf("The amount of days to get notified can not be 0 and larger than %d days", maxNotifyUpfront)
		flashWithData(ctx, errors)
		ctx.Redirect(http.StatusFound, "/account")
		return
	}
	user := getAuthenticatedUser(ctx)
	account, err := db.Store.Account.GetUserAccount(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	account.NotifyUpfront = params.NotifyUpfront
	if err := db.Store.Account.UpdateAccount(account); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	ctx.Redirect(http.StatusFound, "/account")
}

func HandleAccountShow(ctx *gin.Context) {
	user := getAuthenticatedUser(ctx)
	account, err := db.Store.Account.GetUserAccount(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	context := util.Map{
		"account": account,
		"user":    user,
	}
	ctx.HTML(http.StatusOK, "account/show", context)
}
