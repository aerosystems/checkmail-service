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
	ErrReadRequestBody      = ApiError{Message: "could not read request body", HttpCode: http.StatusUnprocessableEntity}
	ErrInvalidDomain        = ApiError{Message: "invalid domain name", HttpCode: http.StatusBadRequest}
	ErrDomainNotFound       = ApiError{Message: "domain not found", HttpCode: http.StatusNotFound}
	ErrDomainInternalCreate = ApiError{Message: "could not create domain", HttpCode: http.StatusInternalServerError}
	ErrDomainInternalGet    = ApiError{Message: "could not find domain", HttpCode: http.StatusInternalServerError}
	ErrDomainInternalUpdate = ApiError{Message: "could not update domain", HttpCode: http.StatusInternalServerError}
	ErrDomainInternalDelete = ApiError{Message: "could not delete domain", HttpCode: http.StatusInternalServerError}
	ErrDomainInternalCount  = ApiError{Message: "could not count domains", HttpCode: http.StatusInternalServerError}
)
