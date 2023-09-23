package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/handlers"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/ssl"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func main() {
	app, err := initApp()
	if err != nil {
		log.Fatal(err)
	}
	db.Init()
	logger.Init()

	ssl.StartCron()

	port := util.GetEnv("LISTEN_PORT", ":3000")
	logger.Log("msg", fmt.Sprintf("Server is listening on port %s...", port))
	log.Fatal(app.Run(port))
}

func initApp() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	gob.Register(util.Map{})
	gob.Register(data.User{})

	store := cookie.NewStore([]byte(util.GetEnv("SESSION_KEY", "secret")))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	})

	scopes := strings.Split(util.GetEnv("GOOGLE_AUTH_SCOPES", ""), ",")
	gothic.Store = store
	goth.UseProviders(
		google.New(
			util.GetEnv("GOOGLE_CLIENT_ID", ""),
			util.GetEnv("GOOGLE_CLIENT_SECRET", ""),
			util.GetEnv("GOOGLE_AUTH_CALLBACK", ""),
			scopes...,
		),
	)

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.New()
	router.Static("/src/static/assets", "./src/static/assets")
	router.HTMLRender = createEngine()

	router.Use(sessions.Sessions("mysession", store))
	router.Use(func(ctx *gin.Context) {
		ctx.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		ctx.Set("Pragma", "no-cache")
		ctx.Set("Expires", "0")
		ctx.Set("Surrogate-Control", "no-store")
		ctx.Next()
	})
	router.Use(handlers.DefaultErrorHandler())
	router.Use(handlers.WithFlash)
	router.Use(handlers.WithViewHelpers)

	router.GET("/", handlers.HandleGetHome)
	//router.GET("/pricing", handlers.HandleGetHome)
	router.GET("/signin", handlers.HandleGetSignin)
	router.POST("/signin", handlers.HandleSigninWithEmail)
	router.POST("/signin/google", handlers.HandleSigninWithGoogle)
	router.GET("/signup", handlers.HandleGetSignup)
	router.POST("/signup", handlers.HandleSignupWithEmail)
	router.GET("/signout", handlers.HandleGetSignout)
	router.GET("/auth/callback", handlers.HandleAuthCallback)

	domains := router.Group("/domains", handlers.WithMustBeAuthenticated)
	domains.GET("/", handlers.HandleDomainList)
	domains.POST("/", handlers.HandleDomainCreate)
	domains.GET("/new", handlers.HandleDomainNew)
	domains.GET("/:id", handlers.HandleDomainShow)
	domains.POST("/:id/delete", handlers.HandleDomainDelete)

	account := router.Group("/account", handlers.WithMustBeAuthenticated)
	account.GET("/", handlers.HandleAccountShow)
	account.POST("/", handlers.HandleAccountUpdate)

	return router, nil
}

func createEngine() *ginview.ViewEngine {
	engine := goview.New(goview.Config{
		Root:      "src/static/views",
		Extension: ".html",
		Funcs: template.FuncMap{
			"css": func(name string) (res template.HTML) {
				filepath.Walk("./src/static/assets", func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.Name() == name {
						res = template.HTML("<link rel=\"stylesheet\" href=\"/" + path + "\">")
					}
					return nil
				})
				return
			},
			"formatTime": func(t time.Time) string {
				timeZero := time.Time{}
				if t.Equal(timeZero) {
					return "n/a"
				}
				return t.Format(time.DateTime)
			},
			"daysLeft": func(t time.Time) string {
				timeZero := time.Time{}
				if t.Equal(timeZero) {
					return "n/a"
				}
				return fmt.Sprintf("%d days", time.Until(t)/(time.Hour*24))
			},
			"badgeForStatus": func(status string) template.HTML {
				switch status {
				case data.StatusHealthy:
					return template.HTML(fmt.Sprintf(`<div class="badge badge-success">%s</div>`, status))
				case data.StatusExpires:
					return template.HTML(fmt.Sprintf(`<div class="badge badge-warning">%s</div>`, status))
				case data.StatusExpired:
					return template.HTML(fmt.Sprintf(`<div class="badge badge-error">%s</div>`, status))
				case data.StatusInvalid:
					return template.HTML(fmt.Sprintf(`<div class="badge badge-error">%s</div>`, status))
				case data.StatusOffline:
					return template.HTML(fmt.Sprintf(`<div class="badge badge-error">%s</div>`, status))
				default:
					return ""
				}
			},
			"pluralize": func(word string, length int) string {
				if length > 1 {
					return fmt.Sprintf("%ss", word)
				}
				return word
			},
			"inc": func(i int) int {
				return i + 1
			},
		},
		DisableCache: true,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})

	return ginview.Wrap(engine)
}
