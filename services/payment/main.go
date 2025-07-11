package main

import (
	"context"
	"letspay/pkg/db"
	"letspay/pkg/rabbitmq"
	"letspay/services/payment/common/constants"
	"letspay/services/payment/config"
	"letspay/services/payment/controller/api"
	"letspay/services/payment/mq"
	"letspay/services/payment/repository/database"
	"letspay/services/payment/repository/provider"
	"letspay/services/payment/repository/provider/midtrans"
	"letspay/services/payment/repository/provider/xendit"
	"letspay/services/payment/scheduler"
	"letspay/services/payment/usecase"
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

	rds := db.InitRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	db := db.InitDB()
	defer db.Close()

	mqConn := rabbitmq.Connect(cfg.RabbitMqUrl)

	disbursementRepo := database.NewDisbursementRepo(db)
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

	disbursementUC := usecase.NewDisbursementUsecase(
		disbursementRepo,
		providerRepo,
		bankRepo,
		rds,
		mqConn,
	)

	scheduler := scheduler.NewScheduler(disbursementUC)
	scheduler.RegisterJobs()

	router := api.HandleRequests(cfg, disbursementUC)

	server := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Server.Port, // set to 0.0.0.0 so docker can listen
		Handler: router,
	}

	var wg sync.WaitGroup

	ctx1, cancel1 := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		// TODO: make init func?
		mq.ConsumeDisbursementRequest(mqConn, disbursementUC.HandleDisbursementRequest)
		ctx1.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("API listening on port: " + cfg.Server.Port)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	ctx2, cancel2 := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		scheduler.Start()
		log.Println("Cron scheduler started")
		<-ctx2.Done() // wait for shutdown signal
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done // when ctrl+c is called signal will be sent here
	log.Println("Shutting down gracefully...")

	cancel1()
	cancel2()

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
