package model

type (
	// TODO: add error

	AppConfig struct {
		Provider map[int]Provider
	}

	Provider struct {
		ClientId     string
		ClientSecret string
		BaseUrl      string
	}
)
