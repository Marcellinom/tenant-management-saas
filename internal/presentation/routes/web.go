package routes

import (
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterRoutes deprecated
func RegisterRoutes(app *provider.Application) {
	route := app.Engine()

	//route.Use(sessions.Sessions("auth-session", cookie.NewStore([]byte("secret"))))
	//
	//route.Static("/static", "public/static")
	//route.LoadHTMLGlob("public/web/*")
	//
	//route.GET("/", func(ctx *gin.Context) {
	//	iam_url := app.Auth().GetProvider() + "/login"
	//	ctx.HTML(http.StatusOK, "index.html", map[string]any{"iam_url": iam_url})
	//})
	route.GET("/healthcheck", func(context *gin.Context) {
		context.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
	registerApis(app)
}
