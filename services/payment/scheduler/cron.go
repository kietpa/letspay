package scheduler

import (
	"context"
	"letspay/services/payment/scheduler/jobs"
	"letspay/services/payment/usecase"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron                *cron.Cron
	disbursementUsecase usecase.DisbursementUsecase
}

func NewScheduler(
	disbursementUsecase usecase.DisbursementUsecase,
) *Scheduler {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger), // panic
	))

	return &Scheduler{
		cron:                c,
		disbursementUsecase: disbursementUsecase,
	}
}

func (s *Scheduler) RegisterJobs() {
	disb := jobs.NewCheckPendingDisbursementsJob(s.disbursementUsecase)
	s.cron.AddJob("*/5 * * * *", disb)
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() context.Context {
	return s.cron.Stop()
}
