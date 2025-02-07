package main

import (
	"backend/internal/core/services"
	"backend/internal/grpc_server/shortlinks"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"
)

func NewShortlinksGrpc(log logger.Logger, shortlinkService services.ShortlinkService, host string) *ShortlinksGrpc {
	return &ShortlinksGrpc{
		handler: NewCreateHandler(log, shortlinkService, host),
	}
}

type ShortlinksGrpc struct {
	shortlinks.UnimplementedShortlinksServer
	handler httpserver.Handler[shortlinkCreateInput, shortlinkCreateOutput]
}

func (s *ShortlinksGrpc) Create(ctx context.Context, req *shortlinks.CreateRequest) (*shortlinks.CreateResponse, error) {
	output, err := s.handler(ctx, shortlinkCreateInput{req.Url})
	if err != nil {
		return nil, err
	}

	return &shortlinks.CreateResponse{
		Link: output.Link,
	}, nil
}
