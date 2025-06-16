package usecase

import (
	"context"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"letspay/repository/provider"
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
	log.Println(err)
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
func (u disbursementUsecase) CreateDisbursement(
	ctx context.Context, createDisbursementRequest model.CreateDisbursementRequest,
) (model.DisbursementDetail, model.Error) {
	// TODO: add validation for amount, userrefid, bankcode

	input := model.CreateDisbursementInput{
		UserId:            123,      // TODO: retrieve user id
		ReferenceId:       "123123", // TODO: make refid generator
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
		log.Println(err)
		return model.DisbursementDetail{}, model.Error{
			Code:    http.StatusInternalServerError,
			Message: constants.INTERNAL_ERROR_MESSAGE,
		}
	}

	// TODO: execute disbursement to providers

	// TODO: if transaction exceed timeout without final error status,
	// respond to user with pending status
	// TODO: add scheduler to check status to provider for pending transactions

	return model.DisbursementDetail{
		ReferenceId:       input.ReferenceId,
		UserReferenceId:   input.UserReferenceId,
		Status:            input.Status,
		Amount:            input.Amount,
		BankCode:          input.BankCode,
		CreatedAt:         input.CreatedAt,
		BankAccountNumber: input.BankAccountNumber,
		BankAccountName:   input.BankAccountName,
		Description:       input.Description,
		FailureCode:       "",
	}, model.Error{}
}
