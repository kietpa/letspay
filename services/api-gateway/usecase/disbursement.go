package usecase

import (
	"context"
	"fmt"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/pkg/util"
	"letspay/services/api-gateway/common/constants"
	"letspay/services/api-gateway/model"
	"letspay/services/api-gateway/repository"
	"net/http"
	"time"
)

type disbursementUsecase struct {
	userRepo repository.UserRepo
}

func NewDisbursementUsecase(userRepo repository.UserRepo) DisbursementUsecase {
	return &disbursementUsecase{
		userRepo: userRepo,
	}
}

// TODO: maybe make a generic func to send disbursement results?
func (u *disbursementUsecase) HandleDisbursementCompleted(input model.DisbursementCompletedEvent) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	webhookUrl, err := u.userRepo.GetUserWebhook(ctx, input.UserId)
	if err != nil {
		logger.Error(ctx, "[Disbursement] Failed to get user webhook") // can improve later
		return
	}

	req := helper.RequestConfig{
		URL:    webhookUrl,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:           input.DisbursementDetail,
		Timeout:        30 * time.Second,
		ExpectedStatus: http.StatusOK,
	}

	logger.Info(ctx, fmt.Sprintf("[Disbursement] sending disbursement to client refid=%s body=%v",
		input.ReferenceId,
		req.Body,
	))

	resp, statusCode, err := helper.SendRequest(req)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Disbursement] error when sending request=%s statusCode=%d refid=%s",
			err,
			statusCode,
			input.ReferenceId,
		))
		return
	}

	logger.Info(ctx, fmt.Sprintf("[Disbursement] sending disbursement to client SUCCESS refid=%s body=%v",
		input.ReferenceId,
		string(resp),
	))
}

func (u *disbursementUsecase) HandleDisbursementFailed(input model.DisbursementFailedEvent) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.PROCESS_ID, util.GenerateRandomHex())

	webhookUrl, err := u.userRepo.GetUserWebhook(ctx, input.UserId)
	if err != nil {
		logger.Error(ctx, "[Disbursement] Failed to get user webhook") // can improve later
		return
	}

	req := helper.RequestConfig{
		URL:    webhookUrl,
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:           input.DisbursementDetail,
		Timeout:        30 * time.Second,
		ExpectedStatus: http.StatusOK,
	}

	logger.Info(ctx, fmt.Sprintf("[Disbursement] sending disbursement to client refid=%s body=%v",
		input.ReferenceId,
		req.Body,
	))

	resp, statusCode, err := helper.SendRequest(req)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Disbursement] error when sending request=%s statusCode=%d refid=%s",
			err,
			statusCode,
			input.ReferenceId,
		))
		return
	}

	logger.Info(ctx, fmt.Sprintf("[Disbursement] sending disbursement to client SUCCESS refid=%s body=%v",
		input.ReferenceId,
		string(resp),
	))
}
