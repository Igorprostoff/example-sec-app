// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"securewebapp/internal/usecase"
	"securewebapp/pkg/logger"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var csrfMd func(http.Handler) http.Handler

// NewRouter -.
// Swagger spec:
// @title       Secure web app API
// @description Using a login form as a test task
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.UserUseCase, csrfSecret string, cookieSecret string) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	store := cookie.NewStore([]byte(cookieSecret))
	handler.Use(sessions.Sessions("vksession", store))
	csrfMd = csrf.Protect([]byte(csrfSecret),
		csrf.MaxAge(0),
		csrf.Secure(true),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"message": "Forbidden - CSRF token invalid"}`))
		})),
	)
	handler.Use(adapter.Wrap(csrfMd))
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))
	handler.LoadHTMLGlob("/html/*.html")
	handler.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request)})
	})

	// Routers
	h := handler.Group("/v1")
	{
		newUserRoutes(h, t.Repo, l)
	}

}
