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
	if authorizationHeader == "" || len("Bearer ") >= len(authorizationHeader) {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	token := authorizationHeader[len("Bearer "):]
	claims, err := decodeJWT(token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid JWT Token: %w", err))
		return
	}
	ctx.Set("token", claims)
	ctx.Next()
}

func CORSMiddleware(allowed_origins ...string) gin.HandlerFunc {
	var origins string
	if len(origins) > 0 {
		origins = strings.Join(allowed_origins, ",")
	} else {
		origins = "*"
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		origin := c.Request.Header.Get("Origin")
		if len(origins) > 0 {
			for _, v := range allowed_origins {
				if v == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origins)
				}
			}
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
