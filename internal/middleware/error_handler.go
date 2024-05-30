package middleware

import (
	"net/http"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/dto"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/apperror"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context) {

	c.Next()

	for _, err := range c.Errors {
		switch e := err.Err.(type) {
		case *apperror.ErrorType:
			{
				c.AbortWithStatusJSON(e.StatusCode,
					dto.Response{
						Msg:  e.Message,
						Data: nil,
					})
			}
		default:
			{
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					dto.Response{
						Msg:  err.Error(),
						Data: nil,
					})
			}

		}
	}

}
