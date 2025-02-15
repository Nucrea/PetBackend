package server

import (
	"backend/cmd/backend/server/handlers"
	"backend/cmd/backend/server/middleware"
	"backend/cmd/backend/server/utils"
	"backend/internal/core/services"
	"backend/internal/integrations"
	"backend/pkg/logger"
	"context"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	logger    logger.Logger
	ginEngine *gin.Engine
}

type NewServerOpts struct {
	DebugMode        bool
	Logger           logger.Logger
	UserService      services.UserService
	ShortlinkService services.ShortlinkService
	Tracer           trace.Tracer
}

func New(opts NewServerOpts) *Server {
	if !opts.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.ContextWithFallback = true // Use it to allow getting values from c.Request.Context()

	r.Static("/webapp", "./webapp")
	r.GET("/health", handlers.NewDummyHandler())

	prometheus := integrations.NewPrometheus()
	r.Any("/metrics", gin.WrapH(prometheus.GetRequestHandler()))

	r.Use(middleware.NewRecoveryMiddleware(opts.Logger, prometheus, opts.DebugMode))
	r.Use(middleware.NewRequestLogMiddleware(opts.Logger, opts.Tracer, prometheus))
	r.Use(middleware.NewTracingMiddleware(opts.Tracer))

	userGroup := r.Group("/user")
	userGroup.POST("/create", handlers.NewUserCreateHandler(opts.Logger, opts.UserService))
	userGroup.POST("/login", handlers.NewUserLoginHandler(opts.Logger, opts.UserService))

	dummyGroup := r.Group("/dummy")
	{
		dummyGroup.Use(middleware.NewAuthMiddleware(opts.UserService))
		dummyGroup.GET("", handlers.NewDummyHandler())
		dummyGroup.POST("/forgot-password", func(c *gin.Context) {
			user := utils.GetUserFromRequest(c)
			opts.UserService.SendEmailForgotPassword(c, user.Id)
		})
	}

	return &Server{
		logger:    opts.Logger,
		ginEngine: r,
	}
}

func (s *Server) Run(ctx context.Context, port uint16) {
	listenAddr := fmt.Sprintf("0.0.0.0:%d", port)
	s.logger.Log().Msgf("server listening on %s", listenAddr)

	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", listenAddr)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("can not create network listener")
	}

	go func() {
		<-ctx.Done()
		s.logger.Log().Msg("stopping tcp listener...")
		listener.Close()
	}()

	err = s.ginEngine.RunListener(listener)
	if err != nil && err == net.ErrClosed {
		s.logger.Fatal().Err(err).Msg("server stopped with error")
	}
}
