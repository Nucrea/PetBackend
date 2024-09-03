package handlers

import (
	"backend/src/core/services"
	"backend/src/logger"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

type shortlinkCreateOutput struct {
	Link string `json:"link"`
}

func NewShortlinkCreateHandler(logger logger.Logger, shortlinkService services.ShortlinkService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxLogger := logger.WithContext(ctx)

		rawUrl := ctx.Query("url")
		if rawUrl == "" {
			ctxLogger.Error().Msg("url query param missing")
			ctx.AbortWithError(400, fmt.Errorf("url query param missing"))
			return
		}

		u, err := url.Parse(rawUrl)
		if err != nil {
			ctxLogger.Error().Err(err).Msg("error parsing url param")
			ctx.Data(400, "plain/text", []byte(err.Error()))
			return
		}
		u.Scheme = "https"

		linkId, err := shortlinkService.CreateShortlink(ctx, u.String())
		if err != nil {
			ctxLogger.Error().Err(err).Msg("err creating shortlink")
			ctx.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		resultBody, err := json.Marshal(shortlinkCreateOutput{
			Link: "https://nucrea.ru/s/" + linkId,
		})
		if err != nil {
			ctxLogger.Error().Err(err).Msg("err marshalling shortlink")
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "application/json", resultBody)
	}
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
