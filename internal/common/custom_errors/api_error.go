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
	ErrInvalidRequestBody      = ApiError{"Invalid request body", http.StatusUnprocessableEntity, codes.InvalidArgument}
	ErrInvalidRequestPayload   = ApiError{"Invalid request payload", http.StatusBadRequest, codes.InvalidArgument}
	ErrInvalidDomain           = ApiError{"Invalid domain name", http.StatusBadRequest, codes.InvalidArgument}
	ErrDomainNotFound          = ApiError{"Domain not found", http.StatusNotFound, codes.NotFound}
	ErrDomainAlreadyExists     = ApiError{"Domain already exists", http.StatusConflict, codes.AlreadyExists}
)
