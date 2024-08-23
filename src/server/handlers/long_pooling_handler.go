package handlers

import (
	"backend/src/client_notifier"
	"backend/src/logger"
	"backend/src/server/utils"

	"github.com/gin-gonic/gin"
)

func NewLongPoolingHandler(logger logger.Logger, notifier client_notifier.ClientNotifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := utils.GetUserFromRequest(c)
		if user == nil {
			c.Data(403, "plain/text", []byte("Unauthorized"))
			return
		}

		eventChan := notifier.RegisterClient(user.Id)

		select {
		case <-c.Done():
			notifier.UnregisterClient(user.Id)
		case event := <-eventChan:
			c.Data(200, "application/json", event.Data)
		}
	}
}
