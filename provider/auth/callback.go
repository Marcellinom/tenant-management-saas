package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (a *Authenticator) RegisterCallback(engine *gin.Engine) {
	engine.GET("/callback", handleCallback)
}

func handleCallback(ctx *gin.Context) {
	iam_callback := os.Getenv("AUTH_PROVIDER") + "/callback"
	ctx.HTML(http.StatusOK, "callback.html", map[string]any{"iam_callback": string(iam_callback)})
}
