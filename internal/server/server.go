package server

import (
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/handler"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HandlerOpts struct {
	User *handler.UserHandler
}

func SetupRouter(opt *HandlerOpts) *gin.Engine {
	router := gin.Default()
	router.ContextWithFallback = true
	router.Use(middleware.HandleError)

	router.POST("/users", opt.User.CreateUser)

	return router

}
