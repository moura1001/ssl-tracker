package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func HandleGetSignup(ctx *gin.Context) {
	data := util.Map{}
	flash, exists := ctx.Get("flash")
	if exists {
		data, _ = flash.(util.Map)
	}

	ctx.HTML(http.StatusOK, "auth/signup", data)
}

type SignupParams struct {
	Email    string `form:"email"`
	FullName string `form:"fullName"`
	Password string `form:"password"`
}

func (p SignupParams) validate() util.Map {
	data := util.Map{}
	if !util.IsValidEmail(p.Email) {
		data["emailError"] = "Please provide a valid email address"
	}
	if !util.IsValidPassword(p.Password) {
		data["passwordError"] = "Please provide a strong password"
	}
	if len(p.FullName) < 3 {
		data["fullNameError"] = "Please provide your real full name"
	}
	return data
}

func HandleSignupWithEmail(ctx *gin.Context) {
	var params SignupParams
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}
	if errors := params.validate(); len(errors) > 0 {
		errors["email"] = params.Email
		errors["fullName"] = params.FullName
		flashWithData(ctx, errors)
		ctx.Redirect(http.StatusFound, "/signup")
		return
	}

	_, err := db.Store.Account.CreateAccountForUserIfNotExist(&data.User{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		errors := util.Map{
			"appError": fmt.Sprintf("error to create account. Details: %s", err),
		}
		errors["email"] = params.Email
		errors["fullName"] = params.FullName
		flashWithData(ctx, errors)
		ctx.Redirect(http.StatusFound, "/signup")
		return
	}

	logger.Log("msg", "user signup with email", "id", "client.Id")
	ctx.HTML(http.StatusOK, "auth-email-confirmation.html", util.Map{
		"email": params.Email,
	})
}

func HandleGetSignin(ctx *gin.Context) {
	checkoutId := ctx.Query("checkoutId")
	ctx.HTML(http.StatusOK, "auth/signin", util.Map{
		"checkoutId": checkoutId,
	})
}

// TODO: implementation
func HandleSigninWithEmail(ctx *gin.Context) {
	fmt.Println("HandleSigninWithEmail")
}

func HandleSigninWithGoogle(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func HandleGetSignout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Set(localsUserKey, nil)
	gothic.Logout(ctx.Writer, ctx.Request)
	if err := session.Save(); err != nil {
		logger.Log("error", "signout error", "err", fmt.Errorf("error in HandleGetSignout saving session: %s", err))
	}
	ctx.Redirect(http.StatusFound, "/")
}

// This is the main callback that will be triggered after each authentication
func HandleAuthCallback(ctx *gin.Context) {
	client, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		u := getAuthenticatedUser(ctx)
		if u == nil {
			ctx.Error(NewDefaultHttpError(fmt.Errorf("invalid access token")))
			ctx.Abort()
			return
		}
	}

	user := data.User{
		Id:          client.UserID,
		Email:       client.Email,
		AccessToken: client.AccessToken,
		ExpiresAt:   client.ExpiresAt,
	}

	session := sessions.Default(ctx)
	session.Set(localsUserKey, user)

	acc, err := db.Store.Account.CreateAccountForUserIfNotExist(&user)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}

	logger.Log("event", "user signin", "userId", user.Id, "accountId", acc.Id)

	if err := session.Save(); err != nil {
		logger.Log("error", "auth error", "err", fmt.Errorf("error in HandleAuthCallback saving session: %s", err))
	}

	ctx.Redirect(http.StatusFound, "/domains")
}
