package constants

const (
	INVALID int = iota
	XENDIT_PROVIDER_ID
	MIDTRANS_PROVIDER_ID
)

const (
	JSON_BODY       = "json_body"
	REQUEST_HEADERS = "headers"
	PROVIDER        = "provider"
	USER_ID         = "user_id"
	PROCESS_ID      = "ProcessID"
	X_USER_ID       = "X-User-ID"

	STATUS_PENDING   = "PENDING"
	STATUS_COMPLETED = "COMPLETED"
	STATUS_FAILED    = "FAILED"

	X_IDEMPOTENCY_KEY = "X-IDEMPOTENCY-KEY"
	WEBHOOK_ID        = "webhook-id"
)
