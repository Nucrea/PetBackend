package httpserver

import (
	"backend/pkg/logger"
	"context"
	"fmt"
	"net"
)

type serverHttp interface {
	RunListener(l net.Listener) error
}

type Server struct {
	logger logger.Logger
	http   serverHttp
}

type NewServerOpts struct {
	Logger     logger.Logger
	HttpServer serverHttp
}

func New(opts NewServerOpts) *Server {
	return &Server{
		logger: opts.Logger,
		http:   opts.HttpServer,
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

	err = s.http.RunListener(listener)
	if err != nil && err == net.ErrClosed {
		s.logger.Fatal().Err(err).Msg("server stopped with error")
	}
}
