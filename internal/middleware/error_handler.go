package middleware

import (
	"net/http"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/dto"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
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
