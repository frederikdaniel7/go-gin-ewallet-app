package server

import (
	"time"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/handler"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HandlerOpts struct {
	User        *handler.UserHandler
	Transaction *handler.TransactionHandler
}

func SetupRouter(opt *HandlerOpts) *gin.Engine {
	router := gin.Default()
	cors := initiateCORS()
	router.Use(cors)
	router.ContextWithFallback = true

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	router.Use(middleware.RequestId)
	router.Use(middleware.Logger(log))
	router.Use(middleware.HandleError)
	router.POST("/users", opt.User.CreateUser)

	router.POST("/login", opt.User.Login)
	router.POST("/password/forgot", opt.User.ForgotPassword)
	router.PATCH("/password/:token", opt.User.ResetPassword)

	router.Use(middleware.AuthHandler)
	router.GET("/users", opt.User.GetUserDetails)
	router.POST("/transactions/transfer", opt.Transaction.Transfer)
	router.POST("/transactions/topup", opt.Transaction.TopUpBalance)
	router.GET("/transactions", opt.Transaction.GetTransactions)
	return router

}

func initiateCORS() gin.HandlerFunc {

	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept", "Access-Control-Allow-Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

}
