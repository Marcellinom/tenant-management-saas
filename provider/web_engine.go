package provider

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type WebEngine = gin.Engine
type Context = gin.Context
type H gin.H

type ErrorHandler interface {
	Handle(debugMode bool) func(ctx *Context)
}

type WebEngineConfig struct {
	DebugMode          bool
	Environtment       string
	Port               string
	CustomErrorHandler ErrorHandler
}

func NewWebEngineConfig(debugMode bool, environtment string, port string) WebEngineConfig {
	return WebEngineConfig{DebugMode: debugMode, Environtment: environtment, Port: port}
}

func DefaultEngineConfig() WebEngineConfig {
	return WebEngineConfig{
		DebugMode:    os.Getenv("APP_DEBUG") == "true",
		Environtment: os.Getenv("APP_ENV"),
		Port:         os.Getenv("PORT"),
	}
}

func (e *WebEngineConfig) UseCustomErrorHandler(handler ErrorHandler) WebEngineConfig {
	e.CustomErrorHandler = handler
	return *e
}

func SetupWebEngine(cfg WebEngineConfig) (*WebEngine, error) {
	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "" {
				return name
			}
			if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "" {
				return name
			}
			return ""
		})
	}

	r.Use(func(ctx *Context) {
		ctx.Set("request_id", uuid.NewString())
	})

	r.NoRoute(func(ctx *Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, H{
			"code":    0,
			"message": "page_not_found",
		})
	})
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(ctx *Context) {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, H{
			"code":    http.StatusMethodNotAllowed,
			"message": "method_not_allowed",
			"data":    nil,
		})
	})
	r.Use(gin.CustomRecovery(func(ctx *Context, err any) {
		requestId, exists := ctx.Get("request_id")
		data := map[string]interface{}{
			"error": "server unable to handle error",
		}
		if exists {
			data["request_id"] = requestId
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, H{
			"code":    500,
			"message": "internal_server_error",
			"data":    data,
		})
	}))
	if cfg.CustomErrorHandler != nil {
		r.Use(cfg.CustomErrorHandler.Handle(cfg.DebugMode))
	}

	log.Println("Gin server successfully set up!")
	return r, nil
}
