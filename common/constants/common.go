package constants

const (
	INVALID int = iota
	XENDIT_PROVIDER_ID
)

const (
	JSON_BODY       = "json_body"
	REQUEST_HEADERS = "headers"
	PROVIDER        = "provider"
	USER_ID         = "user_id"

	STATUS_PENDING   = "PENDING"
	STATUS_COMPLETED = "COMPLETED"
	STATUS_FAILED    = "FAILED"

	// xendit constants
	X_CALLBACK_TOKEN = "x-callback-token"
	WEBHOOK_ID       = "webhook-id"
)
