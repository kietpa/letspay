package model

type (
	// TODO: add error

	AppConfig struct {
		Server      Server
		Provider    map[int]Provider
		Redis       Redis
		RabbitMqUrl string
	}

	Provider struct {
		ClientId      string
		ApiKey        string
		BaseUrl       string
		CallbackToken string
	}

	Server struct {
		Port    string
		Timeout int
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}
)
