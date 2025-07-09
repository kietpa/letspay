package model

type (
	// TODO: add error

	AppConfig struct {
		Server Server
		Redis  Redis
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
