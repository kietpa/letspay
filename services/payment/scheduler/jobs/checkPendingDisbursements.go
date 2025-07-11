package jobs

import (
	"context"
	"fmt"
	"letspay/pkg/logger"
	"letspay/pkg/util"
	"letspay/services/payment/common/constants"
	"letspay/services/payment/usecase"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	logger.Info(ctx, "[Disbursement Scheduler] starting disbursement scheduler...")
	updateCount, err := j.disbursementUsecase.CheckAndUpdatePendingDisbursements(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Disbursement Scheduler] failed to run scheduler err=%s", err))
		return
	}
	logger.Info(ctx, fmt.Sprintf("[Disbursement Scheduler] disbursements updated=%d", updateCount))

}
