package utils

import (
	"backend/src/core/models"

	"github.com/gin-gonic/gin"
)

func GetUserFromRequest(c *gin.Context) *models.UserDTO {
	if user, ok := c.Get("user"); ok {
		return user.(*models.UserDTO)
	}
	return nil
}
