package server

import (
	"backend/cmd/backend/server/handlers"
	"backend/cmd/backend/server/middleware"
	"backend/cmd/backend/server/utils"
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/internal/integrations"
	"backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type NewServerOpts struct {
	DebugMode        bool
	Logger           logger.Logger
	UserService      services.UserService
	ShortlinkService services.ShortlinkService
	Tracer           trace.Tracer
}

func NewServer(opts NewServerOpts) *httpserver.Server {
	if !opts.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.ContextWithFallback = true // Use it to allow getting values from c.Request.Context()

	// r.Static("/webapp", "./webapp")
	r.GET("/health", handlers.NewDummyHandler())

	prometheus := integrations.NewPrometheus()
	r.Any("/metrics", gin.WrapH(prometheus.GetRequestHandler()))

	r.Use(httpserver.NewRecoveryMiddleware(opts.Logger, prometheus, opts.DebugMode))
	r.Use(httpserver.NewRequestLogMiddleware(opts.Logger, opts.Tracer, prometheus))
	r.Use(httpserver.NewTracingMiddleware(opts.Tracer))

	v1 := r.Group("/v1")

	userGroup := v1.Group("/user")
	{
		userGroup.POST("/create", handlers.NewUserCreateHandler(opts.Logger, opts.UserService))
		userGroup.POST("/login", handlers.NewUserLoginHandler(opts.Logger, opts.UserService))
	}

	dummyGroup := v1.Group("/dummy")
	{
		dummyGroup.Use(middleware.NewAuthMiddleware(opts.UserService))
		dummyGroup.GET("", handlers.NewDummyHandler())
		dummyGroup.POST("/forgot-password", func(c *gin.Context) {
			user := utils.GetUserFromRequest(c)
			opts.UserService.ForgotPassword(c, user.Id)
		})
	}

	return httpserver.New(
		httpserver.NewServerOpts{
			Logger:     opts.Logger,
			HttpServer: r,
		},
	)
}
