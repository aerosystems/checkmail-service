package CustomErrors

var apiErrors = []ApiError{
	ErrApiKeyNotFound,
	ErrSubscriptionIsNotActive,
	ErrReadRequestBody,
	ErrInvalidDomain,
	ErrDomainNotFound,
}
