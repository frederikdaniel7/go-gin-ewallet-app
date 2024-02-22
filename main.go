package main

import (
	"log"

	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/handler"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/repository"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/server"
	"git.garena.com/sea-labs-id/bootcamp/batch-03/frederik-hutabarat/assignment-go-rest-api/internal/usecase"
)

func main() {
	if err := ConfigInit(); err != nil {
		log.Fatalf("failed loading env: %s", err.Error())
	}
	InitDB()

	userRepository := repository.NewUserRepository(db)
	walletRepository := repository.NewWalletRepository(db)

	userUseCase := usecase.NewUserUseCaseImpl(userRepository, walletRepository)

	userHandler := handler.NewUserHandler(userUseCase)

	router := server.SetupRouter(&server.HandlerOpts{
		User: userHandler,
	})
	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}