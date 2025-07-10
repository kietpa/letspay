package main

import (
	"context"
	"letspay/pkg/auth"
	"letspay/pkg/rabbitmq"
	"letspay/services/api-gateway/handler"
	"letspay/services/api-gateway/mq"
	"letspay/services/api-gateway/repository"
	"letspay/services/api-gateway/routing"
	"letspay/services/api-gateway/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// TODO: create config to load env and urls
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load env variables, err=%v", err)
		panic(err)
	}

	auth.SetSecret(os.Getenv("JWT_SECRET"))
	mqConn := rabbitmq.Connect(os.Getenv("RABBITMQ_URL"))
	val := validator.New()

	// for internal http calls we use one client
	// TODO: replace with gRPC later?
	httpclient := &http.Client{
		Timeout: 10 * time.Second,
	}

	userRepo := repository.NewUserRepo(
		os.Getenv("USER_URL"),
		httpclient,
	)

	disbursementUC := usecase.NewDisbursementUsecase(userRepo)
	mq.InitConsumers(mqConn, disbursementUC)

	// handler sends messages to queues
	handler := handler.NewApiHandler(mqConn, val)

	router := mux.NewRouter().StrictSlash(true)
	routing.InitRouting(
		router,
		os.Getenv("USER_URL"),
		os.Getenv("PAYMENT_URL"),
		handler,
	)

	server := &http.Server{
		Addr:         "0.0.0.0:8080", // set to 0.0.0.0 so docker can listen
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("API listening on port: 8080")

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	// done blocks untl a value is sent to done
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done // when ctrl+c is called signal will be sent here
	log.Println("Shutting down gracefully...")

	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer httpCancel()
	if err := server.Shutdown(httpShutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("Shutdown complete")
}
