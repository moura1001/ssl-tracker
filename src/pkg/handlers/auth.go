package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func HandleGetSignup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "auth-signup.html", util.Map{})
}

type SignupParams struct {
	Email    string
	FullName string
	Password string
}

func (p SignupParams) Validate() util.Map {
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
	if errors := params.Validate(); len(errors) > 0 {
		errors["email"] = params.Email
		errors["fullName"] = params.FullName
		flashWithData(ctx, errors)
		ctx.Redirect(http.StatusFound, "/signup")
		return
	}

	//client := createSupabaseClient()
	client := createClient()
	/*
		// supabase SignInWithProvider returns a URL for signing in via OAuth
		user, err := client.Auth.User(context.Background(), supabase.UserCredentials {
			Email: params.Email,
			Password: params.Password,
			Data: util.Map{"fullName": params.FullName},
		})
		if err != nil {
			ctx.Error(NewDefaultHttpError(err))
			return
		}

		ourUser := &data.User{Id: user.Id, Email: user.Email}
	*/

	logger.Log("msg", "user signup with email", "id", client.Id)
	ctx.HTML(http.StatusOK, "auth-email-confirmation.html", util.Map{
		"email": params.Email,
	})
}

func HandleGetSignin(ctx *gin.Context) {
	checkoutId := ctx.Query("checkoutId")
	ctx.HTML(http.StatusOK, "auth-signin.html", util.Map{
		"checkoutId": checkoutId,
	})
}

// TODO: implementation
func HandleSigninWithEmail(ctx *gin.Context) {

}

// TODO: implementation
func HandleSigninWithGoogle(ctx *gin.Context) {

}

// TODO: implementation
func HandleGetSignout(ctx *gin.Context) {
	//client := createSupabaseClient()
	//cookie, _ = ctx.Cookie("accessToken")
	//if err := client.Auth.Signout(context.Background(), cookie); err != nil {
	//	ctx.Error(NewDefaultHttpError(fmt.Errorf("invalid access token")))
	//	return
	//}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "accessToken",
		Expires: time.Now().AddDate(0, 0, -10),
		Value:   "",
	})
	ctx.Redirect(http.StatusFound, "/")
}

// This is the main callback that will be triggered after each authentication
func HandleAuthCallback(ctx *gin.Context) {
	accessToken := ctx.Param("accessToken")
	if len(accessToken) <= 0 {
		ctx.Error(NewDefaultHttpError(fmt.Errorf("invalid access token")))
		return
	}
	http.SetCookie(ctx.Writer, &http.Cookie{
		Secure:   false,
		HttpOnly: true,
		Name:     "accessToken",
		Value:    accessToken,
	})
	//ctx.SetCookie("accessToken", accessToken, 100, "/", ctx.Request.Host, false, true)

	//client := createSupabaseClient()
	client := createClient()
	/*
		// supabase SignInWithProvider returns a URL for signing in via OAuth
		user, err := client.Auth.User(context.Background(), supabase.UserCredentials {
			Email: params.Email,
			Password: params.Password,
			Data: util.Map{"fullName": params.FullName},
		})
		if err != nil {
			ctx.Error(NewDefaultHttpError(err))
			return
		}

		ourUser := &data.User{Id: user.Id, Email: user.Email}
	*/

	acc, err := db.Store.Account.CreateAccountForUserIfNotExist(client)
	if err != nil {
		ctx.Error(NewDefaultHttpError(err))
		return
	}

	logger.Log("event", "user signin", "userId", client.Id, acc.Id)

	// check if there is a cookie set with a checkout session to redirect the user to
	// when authenticated
	var (
		checkoutSessionId, _ = ctx.Cookie("checkoutSessionId")
		redirectTo           = "/domains"
	)
	if len(checkoutSessionId) > 0 {
		s := sessions.Default(ctx)
		session := s.Get(checkoutSessionId)
		if session == nil {
			ctx.Error(NewDefaultHttpError(fmt.Errorf("session %s does not exist", checkoutSessionId)))
			return
		}
		// valid session
		if time.Until(time.Unix(0, 0)) > 0 {
			redirectTo = ""
		}
		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:    "checkoutSessionId",
			Expires: time.Now().AddDate(0, 0, -10),
			Value:   "deleted",
		})
	}

	ctx.Redirect(http.StatusFound, redirectTo)
}
