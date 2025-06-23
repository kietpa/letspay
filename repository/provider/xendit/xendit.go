package xendit

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/common/constants"
	"letspay/model"
	"letspay/repository/provider"
	"letspay/tool/helper"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	X_IDEMPOTENCY_KEY = "X-IDEMPOTENCY-KEY"
)

var XenditErr = map[string]string{
	"RECIPIENT_ACCOUNT_NUMBER_ERROR": constants.INVALID_BANK_ACCOUNT_MESSAGE,
	"INVALID_DESTINATION":            constants.INVALID_BANK_ACCOUNT_MESSAGE,
	"TRANSFER_ERROR":                 constants.INVALID_TRANSFER_REQUEST,
	"REJECTED_BY_BANK":               constants.REJECTED_TRANSFER,
	"REJECTED_BY_CHANNEL":            constants.REJECTED_TRANSFER,
}

type (
	xenditExecDisburseInput struct {
		ExternalId        string `json:"external_id"`
		Amount            uint32 `json:"amount"`
		BankCode          string `json:"bank_code"`
		AccountHolderName string `json:"account_holder_name"`
		AccountNumber     string `json:"account_number"`
		Description       string `json:"description"`
	}

	xenditDisbursementObject struct {
		Id                      string `json:"id"`
		ExternalId              string `json:"external_id"`
		UserId                  string `json:"user_id"`
		BankCode                string `json:"bank_code"`
		AccountHolderName       string `json:"account_holder_name"`
		Amount                  uint32 `json:"amount"`
		DisbursementDescription string `json:"disbursement_description"`
		Status                  string `json:"status"`
		FailureCode             string `json:"failure_code,omitempty"`
	}
)

type (
	providerRepo struct {
		baseUrl string
		apiKey  string
	}

	NewProviderRepoInput struct {
		BaseUrl string
		ApiKey  string
	}
)

func NewProviderRepo(input NewProviderRepoInput) provider.ProviderRepo {
	return &providerRepo{
		baseUrl: input.BaseUrl,
		apiKey:  input.ApiKey,
	}
}

// TODO: add bank code converter
func (p *providerRepo) ExecuteDisbursement(
	ctx context.Context, input model.CreateDisbursementInput,
) (model.CreateDisbursementProviderOutput, error) {
	cfg := helper.RequestConfig{
		URL:    p.baseUrl + "/disbursements",
		Method: http.MethodPost,
		Headers: map[string]string{
			X_IDEMPOTENCY_KEY: uuid.New().String(),
			"Content-Type":    "application/json",
		},
		Body: xenditExecDisburseInput{
			ExternalId:        input.ReferenceId,
			Amount:            uint32(input.Amount),
			BankCode:          input.BankCode,
			AccountHolderName: input.BankAccountName,
			AccountNumber:     input.BankAccountNumber,
			Description:       input.Description,
		},
		Timeout: time.Duration(30) * time.Second,
		BasicAuth: &helper.BasicAuthConfig{
			Username: p.apiKey,
			Password: "",
		},
		ExpectedStatus: http.StatusOK,
	}

	resp := xenditDisbursementObject{}
	respByte, statusCode, err := helper.SendRequest(cfg)
	if err != nil {
		fmt.Println(statusCode, "exec disb xendit resp:", err)
		return model.CreateDisbursementProviderOutput{}, err
	}

	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		fmt.Println("exec disb xendit unmarshal:", err)
		return model.CreateDisbursementProviderOutput{}, err
	}

	log.Println(resp)

	output := model.CreateDisbursementProviderOutput{
		ReferenceId:         resp.ExternalId,
		ProviderReferenceId: resp.Id,
		Status:              resp.Status,
		Amount:              float64(resp.Amount),
		BankCode:            resp.BankCode,
		BankAccountName:     resp.AccountHolderName,
		Description:         resp.DisbursementDescription,
		FailureCode:         resp.FailureCode,
	}

	if mssg, ok := XenditErr[resp.FailureCode]; ok {
		output.FailureCode = mssg
		output.StatusCode = http.StatusBadRequest
	}

	return output, nil
}

func (p *providerRepo) GetDisbursementStatus(
	ctx context.Context, providerRefid string,
) (model.GetDisbursementProviderResponse, error) {
	cfg := helper.RequestConfig{
		URL:    p.baseUrl + "/disbursements/" + providerRefid,
		Method: http.MethodGet,
		Headers: map[string]string{
			X_IDEMPOTENCY_KEY: uuid.New().String(),
		},
		Timeout: time.Duration(30) * time.Second,
		BasicAuth: &helper.BasicAuthConfig{
			Username: p.apiKey,
			Password: "",
		},
		ExpectedStatus: http.StatusOK,
	}

	resp := xenditDisbursementObject{}
	respByte, statusCode, err := helper.SendRequest(cfg)
	if err != nil {
		fmt.Println(statusCode, "get disb xendit resp:", err)
		return model.GetDisbursementProviderResponse{}, err
	}

	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		fmt.Println(statusCode, "get disb xendit resp:", err)
		return model.GetDisbursementProviderResponse{}, err
	}

	return model.GetDisbursementProviderResponse{
		ReferenceId:         resp.ExternalId,
		ProviderReferenceId: resp.Id,
		Status:              resp.Status,
		FailureReason:       resp.FailureCode,
	}, nil
}
