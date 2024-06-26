package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/handler"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/repository"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/server"
	"github.com/frederikdaniel7/go-gin-ewallet-app/internal/usecase"
	"github.com/frederikdaniel7/go-gin-ewallet-app/pkg/database"
)

func main() {
	if err := ConfigInit(); err != nil {
		log.Fatalf("failed loading env: %s", err.Error())
	}
	InitDB()

	userRepository := repository.NewUserRepository(db)
	walletRepository := repository.NewWalletRepository(db)
	passwordTokenRepository := repository.NewPasswordTokenRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	transactor := database.NewTransaction(db)

	userUseCase := usecase.NewUserUseCaseImpl(
		userRepository, walletRepository, passwordTokenRepository, transactor)
	transactionUseCase := usecase.NewTransactionUseCaseImpl(userRepository, walletRepository, transactionRepository, transactor)

	userHandler := handler.NewUserHandler(userUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	router := server.SetupRouter(&server.HandlerOpts{
		User:        userHandler,
		Transaction: transactionHandler,
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
