package usecase

import (
	"context"
	"errors"
	"fmt"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/tool/logger"
	"letspay/tool/redis"
	"letspay/tool/util"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewDisbursementUsecase(
	disbursementRepo database.DisbursementRepo,
	providerRepo map[int]provider.ProviderRepo,
	bankRepo database.BankRepo,
	redisRepo *redis.RedisClient,
) DisbursementUsecase {
	return &disbursementUsecase{
		disbursementRepo: disbursementRepo,
		providerRepo:     providerRepo,
		bankRepo:         bankRepo,
		redisRepo:        redisRepo,
	}
}

func (u disbursementUsecase) GetDisbursement(
	ctx context.Context, refid string,
) (model.DisbursementDetail, model.Error) {
	resp, err := u.disbursementRepo.GetDisbursement(ctx, refid)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Get Disbursement] repo get disbursement DB error=%s refid=%s",
			err,
			refid,
		))
		return model.DisbursementDetail{}, model.Error{
			Code:    http.StatusNotFound,
			Message: constants.TRANSACTION_NOT_FOUND_MESSAGE,
		}
	}

	if resp.Status == "PENDING" {
		provResp, err := u.providerRepo[resp.ProviderId].GetDisbursementStatus(ctx, resp.ProviderReferenceId)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("[Get Disbursement] provider get status error=%s refid=%s",
				err,
				refid,
			))
			return model.DisbursementDetail{}, model.Error{
				Code:    http.StatusInternalServerError,
				Message: constants.INTERNAL_ERROR_MESSAGE,
			}
		}

		if provResp.Status != resp.Status {
			err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
				ReferenceId:         resp.ReferenceId,
				ProviderId:          resp.ProviderId,
				ProviderReferenceId: resp.ProviderReferenceId,
				Status:              provResp.Status,
				UpdatedAt:           time.Now(),
				FailureCode:         provResp.FailureReason,
			})
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("[Get Disbursement] provider update DB error=%s refid=%s",
					err,
					refid,
				))
				return model.DisbursementDetail{}, model.Error{
					Code:    http.StatusInternalServerError,
					Message: constants.INTERNAL_ERROR_MESSAGE,
				}
			}

			return model.DisbursementDetail{
				ReferenceId:       resp.ReferenceId,
				UserReferenceId:   resp.UserReferenceId,
				Status:            provResp.Status,
				Amount:            resp.Amount,
				BankCode:          resp.BankCode,
				CreatedAt:         resp.CreatedAt,
				BankAccountNumber: resp.BankAccountNumber,
				BankAccountName:   resp.BankAccountName,
				Description:       resp.Description,
				FailureCode:       provResp.FailureReason,
			}, model.Error{}
		}
	}

	return model.DisbursementDetail{
		ReferenceId:       resp.ReferenceId,
		UserReferenceId:   resp.UserReferenceId,
		Status:            resp.Status,
		Amount:            resp.Amount,
		BankCode:          resp.BankCode,
		CreatedAt:         resp.CreatedAt,
		BankAccountNumber: resp.BankAccountNumber,
		BankAccountName:   resp.BankAccountName,
		Description:       resp.Description,
		FailureCode:       resp.FailureCode,
	}, model.Error{}
}

// flow:
// 1. create disbursement record in DB
// 2. execute disbursement to providers
// 3. record final status, respond to user
// 4. if no final status available, set to pending/hanging
func (u disbursementUsecase) CreateDisbursement(
	ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest, userId int,
) (model.DisbursementDetail, model.Error) {
	input := model.CreateDisbursementInput{
		UserId:            userId,
		ReferenceId:       util.GenerateReferenceId(),
		UserReferenceId:   createDisbursementRequest.UserReferenceId,
		Status:            constants.STATUS_PENDING,
		Amount:            createDisbursementRequest.Amount,
		BankCode:          createDisbursementRequest.BankCode,
		CreatedAt:         time.Now(),
		BankAccountNumber: createDisbursementRequest.BankAccountNumber,
		BankAccountName:   createDisbursementRequest.BankAccountName,
		Description:       createDisbursementRequest.Description,
	}

	logger.Info(ctx, fmt.Sprintf(
		"[Create Disbursement] Creating disbursement refid=%s body=%+v",
		input.ReferenceId,
		input,
	))

	err := u.disbursementRepo.CreateDisbursement(ctx, input)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Create Disbursement] repo insert DB error=%s refid=%s",
			err,
			input.ReferenceId,
		))
		return model.DisbursementDetail{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	// get provider sequence
	bank, err := u.bankRepo.GetBankByCode(ctx, input.BankCode)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Create Disbursement] get bank provider sequence DB error=%s refid=%s",
			err,
			input.ReferenceId,
		))
		if errors.Is(err, pgx.ErrNoRows) {
			return model.DisbursementDetail{}, model.Error{
				Code:    http.StatusBadRequest,
				Message: constants.INVALID_BANK_CODE_MESSAGE,
			}
		} else {
			return model.DisbursementDetail{}, model.Error{
				Code:    http.StatusInternalServerError,
				Message: constants.INTERNAL_ERROR_MESSAGE,
			}
		}
	}

	provSeq := strings.Split(bank.Providers, ",")

	res := model.DisbursementDetail{
		ReferenceId:       input.ReferenceId,
		UserReferenceId:   input.UserReferenceId,
		Status:            "",
		Amount:            input.Amount,
		BankCode:          input.BankCode,
		CreatedAt:         input.CreatedAt,
		BankAccountNumber: input.BankAccountNumber,
		BankAccountName:   input.BankAccountName,
		Description:       input.Description,
		FailureCode:       "",
	}

	for i, provider := range provSeq {
		providerId, _ := strconv.Atoi(strings.TrimSpace(provider))

		resp, err := u.providerRepo[providerId].ExecuteDisbursement(ctx, input)
		if err != nil || resp.Status == "FAILED" {
			// TODO: if transaction exceed timeout without final error status,
			// respond to user with pending/stuck/hanging status
			// if errors.Is(err, context.DeadlineExceeded) {
			// }
			logger.Error(ctx, fmt.Sprintf("[Create Disbursement] provider execute disbursement error=%s refid=%s",
				err,
				input.ReferenceId,
			))
			if i == len(provSeq)-1 {
				// check if error is internal or not
				if resp.Status != "" {
					if err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
						ReferenceId:         input.ReferenceId,
						ProviderId:          providerId,
						ProviderReferenceId: resp.ProviderReferenceId,
						Status:              resp.Status,
						UpdatedAt:           time.Now(),
						FailureCode:         resp.FailureCode,
					}); err != nil {
						logger.Error(ctx, fmt.Sprintf("[Create Disbursement] update disbursement DB error=%s refid=%s",
							err,
							input.ReferenceId,
						))
					}

					if resp.StatusCode == http.StatusBadRequest {
						return model.DisbursementDetail{}, model.Error{
							Code:    http.StatusBadRequest,
							Message: resp.FailureCode,
						}
					}
					return model.DisbursementDetail{}, model.Error{
						Code:    http.StatusInternalServerError,
						Message: constants.INTERNAL_ERROR_MESSAGE,
					}
				} else {
					if err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
						ReferenceId:         input.ReferenceId,
						ProviderId:          0,
						ProviderReferenceId: "",
						Status:              "FAILED",
						UpdatedAt:           time.Now(),
						FailureCode:         "INTERNAL ERROR",
					}); err != nil {
						logger.Error(ctx, fmt.Sprintf("[Create Disbursement] update disbursement DB error=%s refid=%s",
							err,
							input.ReferenceId,
						))
					}

					return model.DisbursementDetail{}, model.Error{
						Code:    http.StatusInternalServerError,
						Message: constants.INTERNAL_ERROR_MESSAGE,
					}
				}
			} else {
				// maybe record the last provider?
				continue
			}
		}

		if resp.Status == "COMPLETED" || resp.Status == "PENDING" {
			// update w/final result
			if err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
				ReferenceId:         input.ReferenceId,
				ProviderId:          providerId,
				ProviderReferenceId: resp.ProviderReferenceId,
				Status:              resp.Status,
				UpdatedAt:           time.Now(),
				FailureCode:         resp.FailureCode,
			}); err != nil {
				logger.Error(ctx, fmt.Sprintf("[Create Disbursement] update disbursement DB error=%s refid=%s",
					err,
					input.ReferenceId,
				))
			}
			res.Status = resp.Status
			res.FailureCode = resp.FailureCode
			break
		}
	}

	return res, model.Error{}
}

func (u disbursementUsecase) CallbackDisbursement(
	ctx context.Context, callbackDisbursementRequest model.CallbackDisbursementRequest,
) model.Error {
	resp, err := u.disbursementRepo.GetDisbursement(ctx, callbackDisbursementRequest.ReferenceId)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf(
			"[Callback Disbursement] get disbursement DB error=%s refid=%s",
			err,
			callbackDisbursementRequest.ReferenceId,
		))
		return model.Error{
			Code:    http.StatusNotFound,
			Message: constants.TRANSACTION_NOT_FOUND_MESSAGE,
		}
	}

	if callbackDisbursementRequest.Status != resp.Status {
		logger.Info(ctx, fmt.Sprintf(
			"[Callback Disbursement] Status mismatch providerStatus=%s internalStatus=%s",
			callbackDisbursementRequest.Status,
			resp.Status,
		))
		err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
			ReferenceId:         resp.ReferenceId,
			ProviderId:          resp.ProviderId,
			ProviderReferenceId: resp.ProviderReferenceId,
			Status:              callbackDisbursementRequest.Status,
			UpdatedAt:           time.Now(),
			FailureCode:         callbackDisbursementRequest.FailureCode,
		})
		if err != nil {
			logger.Error(ctx, fmt.Sprintf(
				"[Callback Disbursement] update disbursement DB error=%s refid=%s",
				err,
				callbackDisbursementRequest.ReferenceId,
			))
			return model.Error{
				Code:    http.StatusInternalServerError,
				Message: constants.INTERNAL_ERROR_MESSAGE,
			}
		}

		err = u.redisRepo.Set(ctx, callbackDisbursementRequest.WebhookId, "exists", time.Duration(5*time.Minute))
		if err != nil {
			logger.Error(ctx, fmt.Sprintf(
				"[Callback Disbursement] set key redis error=%s refid=%s",
				err,
				callbackDisbursementRequest.ReferenceId,
			))
		}
	}

	logger.Info(ctx, fmt.Sprintf("[Callback Disbursement] Successfully proccessed disbursement refid=%s", resp.ReferenceId))

	return model.Error{}
}

func (u disbursementUsecase) CallbackValidateToken(
	ctx context.Context, headers http.Header, provider string,
) bool {
	logger.Info(ctx, fmt.Sprintf("[Validate Token] validating token"))
	switch provider {
	case "xendit":
		return u.providerRepo[constants.XENDIT_PROVIDER_ID].ValidateCallbackToken(
			ctx, headers,
		)
	case "midtrans":
		return true // no validation needed
	}
	return false
}

func (u disbursementUsecase) CheckAndUpdatePendingDisbursements(
	ctx context.Context,
) (int, error) {
	logger.Info(ctx, fmt.Sprintf("[Disbursement Scheduler] getting pending disbursements from DB"))
	disbursements, err := u.disbursementRepo.GetPendingDisbursements(ctx)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Disbursement Scheduler] get pending disbursements DB error=%s", err))
		return 0, err
	}

	count := 0
	// TODO: make this concurrent
	for _, d := range disbursements {
		provResp, err := u.providerRepo[d.ProviderId].GetDisbursementStatus(ctx, d.ProviderReferenceId)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf(
				"[Disbursement Scheduler] provider get disbursement status error=%s refid=%s",
				err,
				d.ReferenceId,
			))
			continue
		}

		logger.Info(ctx, fmt.Sprintf(
			"[Disbursement Scheduler] updating disbursement in DB to %s refid=%s",
			provResp.Status,
			d.ReferenceId,
		))
		err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
			ReferenceId:         d.ReferenceId,
			ProviderId:          d.ProviderId,
			ProviderReferenceId: d.ProviderReferenceId,
			Status:              provResp.Status,
			UpdatedAt:           time.Now(),
			FailureCode:         provResp.FailureReason,
		})
		if err != nil {
			logger.Error(ctx, fmt.Sprintf(
				"[Disbursement Scheduler] update disbursement DB error=%s refid=%s",
				err,
				d.ReferenceId,
			))
			continue
		}
		count++
	}

	return count, nil
}
