package main

import (
	"context"
	"letspay/common/constants"
	"letspay/config"
	"letspay/controller/api"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/repository/provider/midtrans"
	"letspay/repository/provider/xendit"
	"letspay/scheduler"
	"letspay/tool/redis"
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
	// to debug, we have to rebuild from scratch with docker rmi letspay-app
	cfg := config.InitConfig()

	rds := redis.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	db := config.InitDB()
	defer db.Close()

	disbursementRepo := database.NewDisbursementRepo(db)
	userRepo := database.NewUserRepo(db)
	bankRepo := database.NewBankRepo(db)
	xenditRepo := xendit.NewProviderRepo(
		xendit.NewProviderRepoInput{
			BaseUrl:       cfg.Provider[constants.XENDIT_PROVIDER_ID].BaseUrl,
			ApiKey:        cfg.Provider[constants.XENDIT_PROVIDER_ID].ApiKey,
			CallbackToken: cfg.Provider[constants.XENDIT_PROVIDER_ID].CallbackToken,
			RedisRepo:     rds,
		},
	)
	midtransRepo := midtrans.NewProviderRepo(
		midtrans.NewProviderRepoInput{
			BaseUrl:   cfg.Provider[constants.MIDTRANS_PROVIDER_ID].BaseUrl,
			ServerKey: cfg.Provider[constants.MIDTRANS_PROVIDER_ID].ApiKey,
			RedisRepo: rds,
		},
	)
	providerRepo := map[int]provider.ProviderRepo{
		constants.XENDIT_PROVIDER_ID:   xenditRepo,
		constants.MIDTRANS_PROVIDER_ID: midtransRepo,
	}

	disbursementUC := usecase.NewDisbursementUsecase(disbursementRepo, providerRepo, bankRepo, rds)
	userUC := usecase.NewUserUsecase(userRepo)

	scheduler := scheduler.NewScheduler(disbursementUC)
	scheduler.RegisterJobs()

	// TODO: mssg queue

	router := api.HandleRequests(cfg, disbursementUC, userUC)
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
		<-ctx.Done() // wait for shutdown signal
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done // when ctrl+c is called signal will be sent here
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
