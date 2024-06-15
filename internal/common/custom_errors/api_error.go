package CustomErrors

import "net/http"

type ApiError struct {
	Message  string
	HttpCode int
}

func (e ApiError) Error() string {
	return e.Message
}

var (
	ErrReadRequestBody = ApiError{Message: "could not read request body", HttpCode: http.StatusUnprocessableEntity}
	ErrDomainNotFound  = ApiError{Message: "domain not found", HttpCode: http.StatusNotFound}
)
