package usecase

import (
	"context"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
	"letspay/tool/util"
	"log"
	"net/http"
	"time"
)

func NewDisbursementUsecase(
	disbursementRepo database.DisbursementRepo,
	providerRepo map[int]provider.ProviderRepo,
) Disbursement {
	return &disbursementUsecase{
		disbursementRepo: disbursementRepo,
		providerRepo:     providerRepo,
	}
}

func (u disbursementUsecase) GetDisbursement(
	ctx context.Context, refid string,
) (model.DisbursementDetail, model.Error) {
	resp, err := u.disbursementRepo.GetDisbursement(ctx, refid)
	log.Println("usecase get disb repo err:", err)
	if err != nil {
		return model.DisbursementDetail{}, model.Error{
			Code:    http.StatusNotFound,
			Message: constants.TRANSACTION_NOT_FOUND_MESSAGE,
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
// 4. if no final status available, set to pending
func (u disbursementUsecase) CreateDisbursement(
	ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest,
) (model.DisbursementDetail, model.Error) {
	// TODO: add validation for amount, userrefid, bankcode

	input := model.CreateDisbursementInput{
		UserId:            123, // TODO: retrieve user id
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

	err := u.disbursementRepo.CreateDisbursement(ctx, input)
	if err != nil {
		log.Println("usecase create disb repo err:", err)
		return model.DisbursementDetail{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	// TODO: add provider sequence

	// execute disbursement to providers
	resp, err := u.providerRepo[constants.XENDIT_PROVIDER_ID].ExecuteDisbursement(ctx, input)
	if err != nil {
		log.Println("usecase create disb provider repo err:", err)
		if resp.Status == "FAILED" {
			err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
				ReferenceId:         input.ReferenceId,
				ProviderId:          constants.XENDIT_PROVIDER_ID,
				ProviderReferenceId: resp.ProviderReferenceId,
				Status:              resp.Status,
				UpdatedAt:           time.Now(),
				FailureCode:         resp.FailureCode,
			})

			// TODO: handle 400 bad request errors from provider
			return model.DisbursementDetail{}, model.Error{
				Code:    http.StatusInternalServerError,
				Message: constants.INTERNAL_ERROR_MESSAGE,
			}
		} else {
			err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
				ReferenceId:         input.ReferenceId,
				ProviderId:          0,
				ProviderReferenceId: "",
				Status:              "FAILED",
				UpdatedAt:           time.Now(),
				FailureCode:         "INTERNAL ERROR",
			})
			return model.DisbursementDetail{}, model.Error{
				Code:    http.StatusInternalServerError,
				Message: constants.INTERNAL_ERROR_MESSAGE,
			}
		}
	}

	err = u.disbursementRepo.UpdateDisbursement(ctx, model.UpdateDisbursementInput{
		ReferenceId:         input.ReferenceId,
		ProviderId:          constants.XENDIT_PROVIDER_ID,
		ProviderReferenceId: resp.ProviderReferenceId,
		Status:              resp.Status,
		UpdatedAt:           time.Now(),
		FailureCode:         resp.FailureCode,
	})

	// TODO: if transaction exceed timeout without final error status,
	// respond to user with pending status
	// TODO: add scheduler to check status to provider for pending transactions

	return model.DisbursementDetail{
		ReferenceId:       input.ReferenceId,
		UserReferenceId:   input.UserReferenceId,
		Status:            resp.Status,
		Amount:            input.Amount,
		BankCode:          input.BankCode,
		CreatedAt:         input.CreatedAt,
		BankAccountNumber: input.BankAccountNumber,
		BankAccountName:   input.BankAccountName,
		Description:       input.Description,
		FailureCode:       resp.FailureCode,
	}, model.Error{}
}
