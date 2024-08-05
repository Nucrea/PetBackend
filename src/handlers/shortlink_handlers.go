package handlers

import (
	"backend/src/services"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

type shortlinkCreateOutput struct {
	Link string `json:"link"`
}

func NewShortlinkCreateHandler(shortlinkService services.ShortlinkService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawUrl := ctx.Query("url")
		if rawUrl == "" {
			ctx.AbortWithError(400, fmt.Errorf("no url param"))
			return
		}

		u, err := url.Parse(rawUrl)
		if err != nil {
			ctx.Data(500, "plain/text", []byte(err.Error()))
			return
		}
		u.Scheme = "https"

		linkId, err := shortlinkService.CreateLink(u.String())
		if err != nil {
			ctx.Data(500, "plain/text", []byte(err.Error()))
			return
		}

		resultBody, err := json.Marshal(shortlinkCreateOutput{
			Link: "https:/nucrea.ru/s/" + linkId,
		})
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "application/json", resultBody)
	}
}

func NewShortlinkResolveHandler(shortlinkService services.ShortlinkService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		linkId := ctx.Param("linkId")

		linkUrl, err := shortlinkService.GetLink(linkId)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Redirect(301, linkUrl)
	}
}
