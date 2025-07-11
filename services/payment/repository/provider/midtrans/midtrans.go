package midtrans

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
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	midtransExecDisburseInput struct {
		BenficiaryName     string `json:"beneficiary_name"`
		BeneficiaryAccount string `json:"beneficiary_account"`
		BeneficiaryBank    string `json:"beneficiary_bank"`
		BeneficiaryEmail   string `json:"beneficiary_email,omitempty"`
		Amount             uint32 `json:"amount"`
		Notes              string `json:"notes"`
	}

	midtransExecDisburseResponse struct {
		Payout []payout `json:"payout"`
	}

	payout struct {
		Status      string `json:"status"` // queued, processed, completed, failed
		ReferenceNo string `json:"reference_no"`
	}

	midtransGetDisburseResponse struct {
		Amount             uint32    `json:"amount"`
		BenficiaryName     string    `json:"beneficiary_name"`
		BeneficiaryAccount string    `json:"beneficiary_account"`
		Bank               string    `json:"bank"`
		ReferenceNo        string    `json:"reference_no"`
		Notes              string    `json:"notes"`
		BeneficiaryEmail   string    `json:"beneficiary_email,omitempty"`
		Status             string    `json:"status"`
		CreatedBy          string    `json:"created_by"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
	}
)

type (
	providerRepo struct {
		baseUrl   string
		serverKey string
		redisRepo *db.RedisClient
	}

	NewProviderRepoInput struct {
		BaseUrl   string
		ServerKey string
		RedisRepo *db.RedisClient
	}
)

func NewProviderRepo(input NewProviderRepoInput) provider.ProviderRepo {
	return &providerRepo{
		baseUrl:   input.BaseUrl,
		serverKey: input.ServerKey,
		redisRepo: input.RedisRepo,
	}
}

func (p *providerRepo) ExecuteDisbursement(
	ctx context.Context, input model.CreateDisbursementInput,
) (model.CreateDisbursementProviderOutput, error) {
	idem := uuid.New().String()

	cfg := helper.RequestConfig{
		URL:    p.baseUrl + "/iris/api/v1/payouts",
		Method: http.MethodPost,
		Headers: map[string]string{
			constants.X_IDEMPOTENCY_KEY: idem,
			"Content-Type":              "application/json",
			"Accept":                    "application/json",
		},
		Body: midtransExecDisburseInput{
			BenficiaryName:     input.BankAccountName,
			BeneficiaryAccount: input.BankAccountNumber,
			BeneficiaryBank:    strings.ToLower(input.BankCode), // midtrans uses lowercase bank-names
			BeneficiaryEmail:   "beneficiary@example.com",       // for now use dummy email
			Amount:             uint32(input.Amount),
			Notes:              input.Description,
		},
		Timeout: time.Duration(30) * time.Second,
		BasicAuth: &helper.BasicAuthConfig{
			Username: p.serverKey,
			Password: "",
		},
		ExpectedStatus: http.StatusOK,
	}

	logger.Info(ctx, fmt.Sprintf("[Create Disbursement - Provider] creating disbursement at midtrans refid=%s requestBody=%v",
		input.ReferenceId,
		cfg.Body,
	))
	resp := midtransExecDisburseResponse{}
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
		ReferenceId:         input.ReferenceId,
		ProviderReferenceId: resp.Payout[0].ReferenceNo,
		Status:              resp.Payout[0].Status,
		Amount:              input.Amount,
		BankCode:            input.BankCode,
		BankAccountName:     input.BankAccountName,
		Description:         input.Description,
	}

	// midtrans uses lowercase
	if output.Status == "queued" || output.Status == "prcoessed" {
		output.Status = "PENDING"
	} else {
		output.Status = strings.ToUpper(output.Status)
	}

	return output, nil
}

func (p *providerRepo) GetDisbursementStatus(
	ctx context.Context, providerRefid string,
) (model.GetDisbursementProviderResponse, error) {
	cfg := helper.RequestConfig{
		URL:    p.baseUrl + "iris/api/v1/payouts/" + providerRefid,
		Method: http.MethodGet,
		Headers: map[string]string{
			constants.X_IDEMPOTENCY_KEY: uuid.New().String(),
		},
		Timeout: time.Duration(30) * time.Second,
		BasicAuth: &helper.BasicAuthConfig{
			Username: p.serverKey,
			Password: "",
		},
		ExpectedStatus: http.StatusOK,
	}
	logger.Info(ctx, fmt.Sprintf("[Get Disbursement - Provider] checking status at midtrans providerRefid=%s",
		providerRefid,
	))

	resp := midtransGetDisburseResponse{}
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

	// midtrans uses lowercase
	if resp.Status == "queued" || resp.Status == "prcoessed" {
		resp.Status = "PENDING"
	} else {
		resp.Status = strings.ToUpper(resp.Status)
	}

	return model.GetDisbursementProviderResponse{
		ProviderReferenceId: resp.ReferenceNo,
		Status:              resp.Status,
		FailureReason:       "", // midtrans doesn't have failure reason
	}, nil
}

func (p *providerRepo) ValidateCallbackToken(
	ctx context.Context, headers http.Header,
) bool {
	return true
}
