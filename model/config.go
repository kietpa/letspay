package model

type (
	// TODO: add error

	AppConfig struct {
		Server   Server
		Provider map[int]Provider
	}

	Provider struct {
		ClientId     string
		ClientSecret string
		BaseUrl      string
	}

	Server struct {
		Port    string
		Timeout int
	}
)
