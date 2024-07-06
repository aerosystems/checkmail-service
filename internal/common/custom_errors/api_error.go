package CustomErrors

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

type ApiError struct {
	Message  string
	HttpCode int
	GrpcCode codes.Code
}

func (e ApiError) Error() string {
	return e.Message
}

var (
	ErrApiKeyNotFound          = ApiError{"Api key not found", http.StatusUnauthorized, codes.Unauthenticated}
	ErrSubscriptionIsNotActive = ApiError{"Subscription is not active", http.StatusForbidden, codes.PermissionDenied}
	ErrReadRequestBody         = ApiError{"Could not read request body", http.StatusUnprocessableEntity, codes.InvalidArgument}
	ErrInvalidDomain           = ApiError{"Invalid domain name", http.StatusBadRequest, codes.InvalidArgument}
	ErrDomainNotFound          = ApiError{"Domain not found", http.StatusNotFound, codes.NotFound}
)
