package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId(c *gin.Context) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Set("request-id", uuid)
	c.Next()
}
