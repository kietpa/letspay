package scheduler

import (
	"context"
	"letspay/scheduler/jobs"
	"letspay/usecase"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Scheduler struct {
	cron                *cron.Cron
	disbursementUsecase usecase.DisbursementUsecase
	logger              zerolog.Logger
}

func NewScheduler(
	disbursementUsecase usecase.DisbursementUsecase,
	logger zerolog.Logger,
) *Scheduler {
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger), // panic
	))

	return &Scheduler{
		cron:                c,
		disbursementUsecase: disbursementUsecase,
		logger:              logger,
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
