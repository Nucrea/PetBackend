package grpcserver

import (
	"backend/pkg/logger"
	"context"
	"fmt"
	"net"
)

type serverGrpc interface {
	Serve(lis net.Listener) error
}

type Server struct {
	logger logger.Logger
	grpc   serverGrpc
}

type NewServerOpts struct {
	Logger     logger.Logger
	GrpcServer serverGrpc
}

func New(opts NewServerOpts) *Server {
	return &Server{
		logger: opts.Logger,
		grpc:   opts.GrpcServer,
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

	err = s.grpc.Serve(listener)
	if err != nil && err == net.ErrClosed {
		s.logger.Fatal().Err(err).Msg("server stopped with error")
	}
}
