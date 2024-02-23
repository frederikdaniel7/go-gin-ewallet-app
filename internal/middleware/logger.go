package middleware

import (
	"errors"
	"time"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()

		requestId, exist := c.Get("request-id")
		if !exist {
			requestId = ""
		}

		entry := log.WithFields(logrus.Fields{
			"request_id":  requestId,
			"latency":     time.Since(start),
			"method":      c.Request.Method,
			"status_code": statusCode,
			"path":        path,
		})

		// we responds with server errors
		if statusCode >= 500 && statusCode <= 599 {
			var appErr *apperror.ErrorType
			for _, err := range c.Errors {
				if errors.As(err, &appErr) {
					entry.WithField("error", appErr).Error("got error")
					entry.WithField("stack", string(appErr.GetStackTrace())).Error("got error")
				}
			}

			return
		}

		entry.Info("request processed")
	}

}
