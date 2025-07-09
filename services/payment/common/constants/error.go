package constants

const (
	INTERNAL_ERROR_MESSAGE = "Internal error"
	VALIDATION_ERROR       = "Validation error"

	TRANSACTION_NOT_FOUND_MESSAGE = "Transaction not found"
	MISSING_AUTH_HEADER_MESSAGE   = "Authorization header missing"

	METHOD_NOT_ALLOWED_MESSAGE = "Method not allowed"

	INVALID_TOKEN_MESSAGE        = "Token invalid"
	INVALID_BANK_ACCOUNT_MESSAGE = "Bank account invalid"
	INVALID_JSON_BODY            = "Request json body invalid"
	INVALID_PASSWORD_MESSAGE     = "Password invalid"
	INVALID_EMAIL_MESSAGE        = "Email invalid"
	INVALID_BANK_CODE_MESSAGE    = "Bank code invalid"
	INVALID_TRANSFER_REQUEST     = "Transfer details are invalid"

	REJECTED_TRANSFER = "Bank has rejected the transaction, please recheck the details or try again later"
)
