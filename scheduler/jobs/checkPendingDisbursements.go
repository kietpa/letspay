package jobs

import (
	"context"
	"letspay/usecase"
	"log"
	"time"
)

type CheckPendingDisbursementsJob struct {
	disbursementUsecase usecase.DisbursementUsecase
}

func NewCheckPendingDisbursementsJob(input usecase.DisbursementUsecase) *CheckPendingDisbursementsJob {
	return &CheckPendingDisbursementsJob{
		disbursementUsecase: input,
	}
}

func (j *CheckPendingDisbursementsJob) Run() {
	// get all refids of pending trx
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	updateCount, err := j.disbursementUsecase.CheckAndUpdatePendingDisbursements(ctx)
	if err != nil {
		log.Println("check pending get disb err: ", err)
		return
	}
	log.Println("disbursements updated: ", updateCount)

}
