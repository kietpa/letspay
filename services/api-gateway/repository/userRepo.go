package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"letspay/pkg/logger"
	"letspay/services/api-gateway/model"
	"net/http"
	"strconv"
)

type (
	UserRepo interface {
		GetUserWebhook(ctx context.Context, userId int) (string, error)
	}

	userRepo struct {
		userServiceUrl string
		httpclient     *http.Client
	}
)

func NewUserRepo(
	userServiceUrl string,
	httpclient *http.Client,
) UserRepo {
	return &userRepo{
		userServiceUrl: userServiceUrl,
		httpclient:     httpclient,
	}
}

func (r *userRepo) GetUserWebhook(ctx context.Context, userId int) (string, error) {
	url := r.userServiceUrl + "/user/internal/" + strconv.Itoa(userId)
	data := model.GetUserDetail{}

	// for debug & info
	logger.Info(ctx, fmt.Sprintf("[Get Webhook] url=%v userid=%d", url, userId))

	resp, err := r.httpclient.Get(url)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("[Get Webhook] http get error=%s", err))
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error(ctx, fmt.Sprintf("[Get Webhook] unmarshal error=%s", err))
		return "", err
	}

	logger.Info(ctx, fmt.Sprintf("[Get Webhook] unmarshalled body=%+v", data))

	return data.Webhook, nil
}
