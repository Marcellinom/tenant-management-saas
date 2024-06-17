package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Authenticator) RegisterCallback(engine *gin.Engine) {
	engine.GET("/callback", a.handleCallback)
}

func (a *Authenticator) handleCallback(ctx *gin.Context) {
	iam_callback := a.auth_provider + "/callback"
	ctx.HTML(http.StatusOK, "callback.html", map[string]any{"iam_callback": iam_callback})
}
