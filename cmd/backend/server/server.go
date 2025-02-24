package server

import (
	"backend/cmd/backend/server/handlers"
	"backend/cmd/backend/server/middleware"
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

	metrics := integrations.NewMetrics("backend")
	serverMetrics := httpserver.NewServerMetrics(metrics)

	r.GET("/health", handlers.New200OkHandler())
	r.Any("/metrics", gin.WrapH(metrics.HttpHandler()))

	r.Use(httpserver.NewRecoveryMiddleware(opts.Logger, serverMetrics, opts.DebugMode))
	r.Use(httpserver.NewRequestLogMiddleware(opts.Logger, opts.Tracer, serverMetrics))
	r.Use(httpserver.NewTracingMiddleware(opts.Tracer))

	r.GET("/verify-user", handlers.NewUserVerifyEmailHandler(opts.Logger, opts.UserService))

	api := r.Group("/api")

	v1 := api.Group("/v1")
	userGroup := v1.Group("/user")
	{
		userGroup.POST("/create", handlers.NewUserCreateHandler(opts.Logger, opts.UserService))
		userGroup.POST("/login", handlers.NewUserLoginHandler(opts.Logger, opts.UserService))
		userGroup.POST("/send-verify", handlers.NewUserSendVerifyEmailHandler(opts.Logger, opts.UserService))
		userGroup.POST("/send-restore-password", handlers.NewUserSendRestorePasswordHandler(opts.Logger, opts.UserService))
		userGroup.POST("/restore-password", handlers.NewUserRestorePasswordHandler(opts.Logger, opts.UserService))

		userGroup.Use(middleware.NewAuthMiddleware(opts.UserService))
		userGroup.POST("/change-password", handlers.NewUserChangePasswordHandler(opts.Logger, opts.UserService))
	}

	return httpserver.New(
		httpserver.NewServerOpts{
			Logger:     opts.Logger,
			HttpServer: r,
		},
	)
}
