package dto

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}

const (
	CodeSessionExpired        = "SESSION_EXPIRED"
	CodeMissingOAuthCode      = "MISSING_OAUTH_CODE"
	CodeMissingTokenID        = "MISSING_TOKEN_ID"
	CodeInvalidGoogleToken    = "INVALID_GOOGLE_TOKEN"
	CodeMissingGoogleEmail    = "MISSING_GOOGLE_EMAIL"
	CodeUnverifiedGoogleEmail = "UNVERIFIED_GOOGLE_EMAIL"
	CodeInternalError         = "INTERNAL_SERVER_ERROR"
	CodeInvalidRequestBody    = "INVALID_REQUEST_BODY"
	CodeMissingUserIdContext  = "MISSING_USER_ID"
	CodeMissingCookie         = "MISSING_COOKIE"
	CodeMissingHeader         = "MISSING_HEADER"
	CodeInvalidHeaderFormat   = "INVALID_HEADER_FORMAT"
	CodeInvalidAccessToken    = "INVALID_ACCESS_TOKEN"
	CodeRestrictedArea        = "RESTRICTED_AREA"
	CodeNoPermission          = "NO_PERMISSION"
	CodePageNotFound          = "PAGE_NOT_FOUND"
	CodeTooManyRequests       = "TOO_MANY_REQUESTS"
)

const (
	MsgSessionExpired         = "Session expired or state security mismatch."
	MsgUnexpectedError        = "An unexpected error occurred."
	MsgMissingRefreshCookie   = "Refresh cookie is required."
	MsgMissingDataFromContext = "Authentication data missing from the context."
	MsgRestrictedArea         = "Restricted area."
	MsgPageNotFound           = "Page not found."
	MsgNoPermission           = "You don't have permission to do it."
	MsgTooManyRequests        = "Too many requests. Please try again later."
)
