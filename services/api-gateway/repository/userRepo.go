package repository

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type (
	UserRepo interface {
		GetUserWebhook(userId int) (string, error)
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

// TODO: create get user webhook function in user service
func (r *userRepo) GetUserWebhook(userId int) (string, error) {
	url := r.userServiceUrl + "/" + strconv.Itoa(userId)
	var webhookUrl string

	resp, err := r.httpclient.Get(url)
	if err != nil {
		return "", err
	}

	body, _ := io.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &webhookUrl); err != nil {
		return "", err
	}

	return "", nil
}
