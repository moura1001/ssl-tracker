package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
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

func WithAuthenticatedUser(ctx *gin.Context) {
	ctx.Set(localsUserKey, nil)
	//client := createSupabaseClient()
	_ = createClient()
	token, err := ctx.Cookie("accessToken")

	if err != nil || len(token) <= 0 {
		ctx.Next()
		return
	}

	/*
		// supabase SignInWithProvider returns a URL for signing in via OAuth
		user, err := client.Auth.User(context.Background(), token)
		if err != nil {
			logger.Log("error", "authentication error", "err", "probably invalid access token")
			http.SetCookie(ctx.Writer, &http.Cookie{
				Name:    "accessToken",
				Expires: time.Now().AddDate(0, 0, -10),
				Value:   "",
			})
			//ctx.SetCookie("accessToken", "", -1, "/", ctx.Request.Host, false, true)
			ctx.Redirect(http.StatusFound, "/")
			return
		}

		ourUser := &data.User{Id: user.Id, Email: user.Email}
	*/
	ourUser := &data.User{Id: "123", Email: "email@email.com"}
	ctx.Set(localsUserKey, ourUser)

	ctx.Next()
}

func WithMustBeAuthenticated(ctx *gin.Context) {
	if !isUserSignedIn(ctx) {
		ctx.Redirect(http.StatusFound, "/")
	}
	ctx.Next()
}

func createClient() *data.User {
	return nil
}
