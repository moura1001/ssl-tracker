package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"time"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/handlers"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

func main() {
	app, err := initApp()
	if err != nil {
		log.Fatal(err)
	}
	logger.Init()

	//ssl.StartCron()

	logger.Log("msg", "Server is listening on port 3000...")
	app.Run(":3000")
}

func initApp() (*gin.Engine, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	store := cookie.NewStore([]byte(util.GetEnv("SESSION_KEY", "secret")))

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	router := gin.New()
	//router.LoadHTMLGlob("src/static/views/**/*.html")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	//config := router.Static("/src/static/assets", "./src/static/assets")
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
	router.Use(handlers.WithAuthenticatedUser)
	router.Use(handlers.WithViewHelpers)

	router.GET("/account", handlers.HandleAccountShow)

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
				return t.Format(time.RFC1123Z)
			},
			"daysLeft": func(t time.Time) string {
				timeZero := time.Time{}
				if t.Equal(timeZero) {
					return "n/a"
				}
				return fmt.Sprintf("%d days", time.Until(t)/(time.Hour*24))
			},
			"badgeForStatus": func(status string) string {
				switch status {
				case data.StatusHealthy:
					return fmt.Sprintf(`<badge class="badge badge-success">%s</badge>`, status)
				case data.StatusExpires:
					return fmt.Sprintf(`<badge class="badge badge-warning">%s</badge>`, status)
				case data.StatusExpired:
					return fmt.Sprintf(`<badge class="badge badge-error">%s</badge>`, status)
				case data.StatusInvalid:
					return fmt.Sprintf(`<badge class="badge badge-error">%s</badge>`, status)
				case data.StatusOffline:
					return fmt.Sprintf(`<badge class="badge badge-error">%s</badge>`, status)
				default:
					return ""
				}
			},
		},
		DisableCache: true,
		Delims:       goview.Delims{Left: "{{", Right: "}}"},
	})

	return ginview.Wrap(engine)
}
