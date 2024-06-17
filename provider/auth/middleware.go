package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated(ctx *gin.Context) {
	authorizationHeader := ctx.Request.Header.Get("Authorization")

	if authorizationHeader == "" {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	token := authorizationHeader[len("Bearer "):]
	claims, err := decodeJWT(token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid JWT Token"))
		return
	}
	ctx.Set("token", claims)
	ctx.Next()
}

func CORSMiddleware(origins ...string) gin.HandlerFunc {
	var allowed_origins string
	if len(origins) > 0 {
		allowed_origins = strings.Join(origins, ",")
	} else {
		allowed_origins = "*"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowed_origins)
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
