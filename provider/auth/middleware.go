package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated(ctx *gin.Context) {
	fmt.Println("session: ", sessions.Default(ctx).Get("profile"))
	fmt.Println(os.Getenv("AUTH_PROVIDER")+"/login", sessions.Default(ctx).Get("profile"))
	res, err := http.Get(os.Getenv("AUTH_PROVIDER") + "/login")
	if err != nil {
		ctx.Error(err)
	}
	defer res.Body.Close()
	payload, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.Error(err)
	}

	fmt.Println(res.StatusCode, string(payload))
	if sessions.Default(ctx).Get("profile") == nil {
		//os.Exit(1)

		//os.Exit(1)
		//var LoginResponse struct {
		//	Url string `json:"url"`
		//}
		//err = json.Unmarshal(payload, &LoginResponse)
		//if err != nil {
		//	ctx.Error(err)
		//}
		//ctx.Redirect(http.StatusSeeOther, LoginResponse.Url)
		ctx.Abort()
	} else {
		ctx.Next()
	}
}

func CORSMiddleware(origins ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
