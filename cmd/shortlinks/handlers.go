package main

import (
	"backend/internal/core/services"
	httpserver "backend/internal/http_server"
	"backend/pkg/logger"
	"context"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

type shortlinkCreateInput struct {
	Url string `json:"url"`
}

type shortlinkCreateOutput struct {
	Link string `json:"link"`
}

func NewShortlinkCreateHandler(
	log logger.Logger,
	shortlinkService services.ShortlinkService,
	serviceUrl string,
) httpserver.Handler[shortlinkCreateInput, shortlinkCreateOutput] {
	return func(ctx context.Context, input shortlinkCreateInput) (shortlinkCreateOutput, error) {
		output := shortlinkCreateOutput{}

		u, err := url.Parse(input.Url)
		if err != nil {
			return output, err
		}
		u.Scheme = "https"

		linkId, err := shortlinkService.CreateShortlink(ctx, u.String())
		if err != nil {
			return output, err
		}

		return shortlinkCreateOutput{
			Link: fmt.Sprintf("%s/s/%s", serviceUrl, linkId),
		}, nil
	}
}

func NewShortlinkCreateGinHandler(
	log logger.Logger,
	shortlinkService services.ShortlinkService,
	serviceUrl string,
) gin.HandlerFunc {
	return httpserver.WrapGin(log, NewShortlinkCreateHandler(log, shortlinkService, serviceUrl))
}

func NewShortlinkResolveHandler(logger logger.Logger, shortlinkService services.ShortlinkService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxLogger := logger.WithContext(ctx)

		linkId := ctx.Param("linkId")

		linkUrl, err := shortlinkService.GetShortlink(ctx, linkId)
		if err == services.ErrShortlinkNotexist {
			ctxLogger.Error().Err(err).Msg("err getting shortlink")
			ctx.AbortWithError(404, err)
			return
		}
		if err == services.ErrShortlinkExpired {
			ctxLogger.Error().Err(err).Msg("err getting shortlink")
			ctx.AbortWithError(404, err)
			return
		}
		if err != nil {
			ctxLogger.Error().Err(err).Msg("unexpected err getting shortlink")
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Redirect(301, linkUrl)
	}
}
