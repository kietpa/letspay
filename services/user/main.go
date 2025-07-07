package main

import (
	"context"
	"letspay/pkg/db"
	"letspay/services/user/config"
	"letspay/services/user/controller/api"
	"letspay/services/user/repository/database"
	"letspay/services/user/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	cfg := config.InitConfig()

	// rds := db.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	db := db.InitDB()
	defer db.Close()

	userRepo := database.NewUserRepo(db)
	userUC := usecase.NewUserUsecase(userRepo)

	// TODO: mssg queue

	router := api.HandleRequests(cfg, userUC)
	server := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Server.Port, // set to 0.0.0.0 so docker can listen
		Handler: router,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("API listening on port: " + cfg.Server.Port)

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
