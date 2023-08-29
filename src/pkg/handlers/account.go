package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
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
	account, err := data.GetUserAccount(user.Id)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	account.NotifyUpfront = params.NotifyUpfront
	if err := data.UpdateAccount(account); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	ctx.Redirect(http.StatusFound, "/account")
}

// TODO: implementation
func HandleAccountShow(ctx *gin.Context) {
	account := data.Account{
		Email:              "email@email.com",
		Plan:               data.PlanFree,
		NotifyUpfront:      4,
		DefaultNotifyEmail: "email@email.com",
	}

	data := util.Map{
		"user":    account,
		"account": account,
	}
	ctx.HTML(http.StatusOK, "account/show", data)
}
