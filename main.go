package main

import (
	"fmt"
	"io"
	"letspay/common/constants"
	"letspay/config"
	"letspay/util"
	"net/http"
)

func main() {
	// config
	cfg := config.InitConfig()

	// logger
	// redis & DB
	db := config.InitDB()

	// init repo with DB instance

	// init providers
	// provider mapper

	// scheduler
	// routing/handler\

	url := cfg.Provider[constants.BRICK_PROVIDER_ID].BaseUrl + "/payments/auth/token"

	req, _ := http.NewRequest("GET", url, nil)

	auth := util.Base64Encode(cfg.Provider[constants.BRICK_PROVIDER_ID].ClientId + ":" + cfg.Provider[constants.BRICK_PROVIDER_ID].ClientSecret)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Basic "+auth)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))

	db.Close()
}
