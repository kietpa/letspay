package main

import (
	"context"
	"letspay/services/api-gateway/routing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// to debug, we have to rebuild from scratch with docker rmi letspay-app
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to load env variables, err=%v", err)
		panic(err)
	}

	// TODO: mssg queue

	router := mux.NewRouter().StrictSlash(true)
	routing.InitRouting(
		router,
		os.Getenv("USER_URL"),
		os.Getenv("PAYMENT_URL"),
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
