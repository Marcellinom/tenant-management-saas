package routes

import (
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// RegisterRoutes deprecated
func RegisterRoutes(app *provider.Application) {
	route := app.Engine()

	route.Use(sessions.Sessions("auth-session", cookie.NewStore([]byte("secret"))))

	route.Static("/static", "public/static")
	route.LoadHTMLGlob("public/web/*")

	route.GET("/", func(ctx *gin.Context) {
		iam_url := os.Getenv("AUTH_PROVIDER") + "/login"
		ctx.HTML(http.StatusOK, "index.html", map[string]any{"iam_url": string(iam_url)})
	})
}
