package main

import (
	"letspay/config"
	"letspay/controller/api"
	"letspay/repository/database"
)

func main() {
	// config
	cfg := config.InitConfig()

	// logger
	// redis & DB
	db := config.InitDB()

	// init repo with DB instance
	disbursementRepo := database.NewDisbursementRepo(db)

	// TODO: init providers (agents in test)
	// provider mapper

	// scheduler
	// mssg queue

	// routing/handler
	api.HandleRequests(cfg, disbursementRepo)

	// url := cfg.Provider[constants.BRICK_PROVIDER_ID].BaseUrl + "/payments/auth/token"

	// req, _ := http.NewRequest("GET", url, nil)

	// auth := util.Base64Encode(cfg.Provider[constants.BRICK_PROVIDER_ID].ClientId + ":" + cfg.Provider[constants.BRICK_PROVIDER_ID].ClientSecret)

	// req.Header.Add("accept", "application/json")
	// req.Header.Add("authorization", "Basic "+auth)

	// res, _ := http.DefaultClient.Do(req)

	// defer res.Body.Close()
	// body, _ := io.ReadAll(res.Body)

	// fmt.Println(string(body))

	db.Close()
}
