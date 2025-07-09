package xendit

import (
	"context"
	"encoding/json"
	"fmt"
	"letspay/pkg/db"
	"letspay/pkg/helper"
	"letspay/pkg/logger"
	"letspay/services/payment/common/constants"
	"letspay/services/payment/model"
	"letspay/services/payment/repository/provider"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	X_CALLBACK_TOKEN = "x-callback-token"
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
		baseUrl       string
		apiKey        string
		callbackToken string
		redisRepo     *db.RedisClient
	}

	NewProviderRepoInput struct {
		BaseUrl       string
		ApiKey        string
		CallbackToken string
		RedisRepo     *db.RedisClient
	}
)

func NewProviderRepo(input NewProviderRepoInput) provider.ProviderRepo {
	return &providerRepo{
		baseUrl:       input.BaseUrl,
		apiKey:        input.ApiKey,
		callbackToken: input.CallbackToken,
		redisRepo:     input.RedisRepo,
	}
}

// TODO: add bank code mapper
func (p *providerRepo) ExecuteDisbursement(
	ctx context.Context, input model.CreateDisbursementInput,
) (model.CreateDisbursementProviderOutput, error) {
	idem := uuid.New().String()

	cfg := helper.RequestConfig{
		URL:    p.baseUrl + "/disbursements",
		Method: http.MethodPost,
		Headers: map[string]string{
			constants.X_IDEMPOTENCY_KEY: idem,
			"Content-Type":              "application/json",
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

	logger.Info(ctx, fmt.Sprintf("[Create Disbursement - Provider] creating disbursement at xendit refid=%s requestBody=%v",
		input.ReferenceId,
		cfg.Body,
	))

	resp := xenditDisbursementObject{}
	respByte, statusCode, err := helper.SendRequest(cfg)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Create Disbursement - Provider] error when sending request=%s statusCode=%d refid=%s",
			err,
			statusCode,
			input.ReferenceId,
		))
		return model.CreateDisbursementProviderOutput{}, err
	}

	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return model.CreateDisbursementProviderOutput{}, err
	}

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
			constants.X_IDEMPOTENCY_KEY: uuid.New().String(),
		},
		Timeout: time.Duration(30) * time.Second,
		BasicAuth: &helper.BasicAuthConfig{
			Username: p.apiKey,
			Password: "",
		},
		ExpectedStatus: http.StatusOK,
	}
	logger.Info(ctx, fmt.Sprintf("[Get Disbursement - Provider] checking status at xendit providerRefid=%s",
		providerRefid,
	))

	resp := xenditDisbursementObject{}
	respByte, statusCode, err := helper.SendRequest(cfg)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Get Disbursement - Provider] error when sending request=%s statusCode=%d providerRefid=%s",
			err,
			statusCode,
			providerRefid,
		))
		return model.GetDisbursementProviderResponse{}, err
	}

	err = json.Unmarshal(respByte, &resp)
	if err != nil {
		return model.GetDisbursementProviderResponse{}, err
	}

	return model.GetDisbursementProviderResponse{
		ProviderReferenceId: resp.Id,
		Status:              resp.Status,
		FailureReason:       resp.FailureCode,
	}, nil
}

func (p *providerRepo) ValidateCallbackToken(
	ctx context.Context, headers http.Header,
) bool {
	xCallbackToken := headers.Get(X_CALLBACK_TOKEN) // in env
	webhookId := headers.Get(constants.WEBHOOK_ID)  // for idempotency
	logger.Info(ctx, fmt.Sprintf("[Validate Token - Provider] validating x-callback-token=%s webhook-id=%s",
		xCallbackToken,
		webhookId,
	))

	return p.callbackToken == xCallbackToken && !p.redisRepo.Exists(ctx, webhookId)
}
