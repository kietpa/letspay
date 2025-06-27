package main

import (
	"context"
	"letspay/common/constants"
	"letspay/config"
	"letspay/controller/api"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/repository/provider/xendit"
	"letspay/scheduler"
	"letspay/tool/logger"
	"letspay/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "letspay/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Letspay
//	@version		1.0
//	@description	Payment Aggregator App

//	@contact.name	Kiet Asmara
//	@contact.url	https://kietpa.github.io/
//	@contact.email	kiet123pascal@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host	localhost:8080
func main() {
	cfg := config.InitConfig()

	// TODO: redis
	db := config.InitDB()
	defer db.Close()

	// logCustom := zlog.With().Timestamp().Logger()

	logCustom := logger.New(
		logger.Config{
			FilePath:   "var/log/app.log", //TODO: tidy this up
			MaxSizeMB:  100,
			MaxBackups: 2,
			MaxAgeDays: 28,
			Compress:   false,
			LokiURL:    "http://loki:3100/loki/api/v1/push",
			LokiLabels: map[string]string{
				"app": "letspay",
			},
		},
	)

	disbursementRepo := database.NewDisbursementRepo(db)
	userRepo := database.NewUserRepo(db)
	xenditRepo := xendit.NewProviderRepo(
		xendit.NewProviderRepoInput{
			BaseUrl:       cfg.Provider[constants.XENDIT_PROVIDER_ID].BaseUrl,
			ApiKey:        cfg.Provider[constants.XENDIT_PROVIDER_ID].ApiKey,
			CallbackToken: cfg.Provider[constants.XENDIT_PROVIDER_ID].CallbackToken,
		},
	)
	providerRepo := map[int]provider.ProviderRepo{
		constants.XENDIT_PROVIDER_ID: xenditRepo,
	}

	disbursementUC := usecase.NewDisbursementUsecase(disbursementRepo, providerRepo, logCustom)
	userUC := usecase.NewUserUsecase(userRepo, logCustom)

	scheduler := scheduler.NewScheduler(disbursementUC, logCustom)
	scheduler.RegisterJobs()

	// TODO: mssg queue

	router := api.HandleRequests(cfg, disbursementUC, userUC, logCustom)

	router.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

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

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start()
		log.Println("Cron scheduler started")
		<-ctx.Done() // Wait for shutdown signal
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	log.Println("Shutting down gracefully...")

	cancel()

	httpShutdownCtx, httpCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer httpCancel()
	if err := server.Shutdown(httpShutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	cronStopCtx := scheduler.Stop()
	select {
	case <-cronStopCtx.Done():
		log.Println("Cron jobs completed")
	case <-time.After(15 * time.Second):
		log.Println("Timeout waiting for cron jobs to complete")
	}

	wg.Wait()
	log.Println("Shutdown complete")

}
