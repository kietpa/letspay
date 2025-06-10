package usecase

import (
	"context"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/database"
	"net/http"
)

func NewDisbursementUsecase(
	disbursementRepo database.DisbursementRepo,
) Disbursement {
	return &disbursementUsecase{
		disbursementRepo: disbursementRepo,
	}
}

func (u disbursementUsecase) GetDisbursement(
	ctx context.Context, refid string,
) (model.DisbursementDetail, model.Error) {
	resp, err := u.disbursementRepo.GetDisbursement(ctx, refid)
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
		CreatedAt:         resp.CreatedAt,
		BankAccountNumber: resp.BankAccountNumber,
		BankAccountName:   resp.BankAccountName,
		Description:       resp.Description,
	}, model.Error{}
}
