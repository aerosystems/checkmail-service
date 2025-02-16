package entities

import (
	"github.com/aerosystems/common-service/customerrors"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	ErrApiKeyNotFound          = customerrors.InternalError{Message: "Api key not found", HttpCode: http.StatusUnauthorized, GrpcCode: codes.Unauthenticated}
	ErrSubscriptionIsNotActive = customerrors.InternalError{Message: "Subscription is not active", HttpCode: http.StatusForbidden, GrpcCode: codes.PermissionDenied}
	ErrInvalidRequestBody      = customerrors.InternalError{Message: "Invalid request body", HttpCode: http.StatusUnprocessableEntity, GrpcCode: codes.InvalidArgument}
	ErrInvalidRequestPayload   = customerrors.InternalError{Message: "Invalid request payload", HttpCode: http.StatusBadRequest, GrpcCode: codes.InvalidArgument}
	ErrInvalidDomain           = customerrors.InternalError{Message: "Invalid domain name", HttpCode: http.StatusBadRequest, GrpcCode: codes.InvalidArgument}
	ErrDomainNotFound          = customerrors.InternalError{Message: "Domain not found", HttpCode: http.StatusNotFound, GrpcCode: codes.NotFound}
	ErrDomainAlreadyExists     = customerrors.InternalError{Message: "Domain already exists", HttpCode: http.StatusConflict, GrpcCode: codes.AlreadyExists}
	ErrInternalError           = customerrors.InternalError{Message: "Internal error", HttpCode: http.StatusInternalServerError, GrpcCode: codes.Internal}
)

var (
	ErrEmailNotValid             = customerrors.ExternalError{Code: 400001, Message: "email address does not valid", HttpCode: http.StatusBadRequest}
	ErrDomainNotValid            = customerrors.ExternalError{Code: 400002, Message: "domain does not valid", HttpCode: http.StatusBadRequest}
	ErrDomainNotExist            = customerrors.ExternalError{Code: 400003, Message: "domain does not exist", HttpCode: http.StatusBadRequest}
	ErrDomainTrustedTypes        = customerrors.ExternalError{Code: 400003, Message: "domain type does not exist in trusted types", HttpCode: http.StatusBadRequest}
	ErrDomainCoverage            = customerrors.ExternalError{Code: 400004, Message: "domain coverage does not exist in trusted coverages", HttpCode: http.StatusBadRequest}
	ErrAccessSubscriptionExpired = customerrors.ExternalError{Code: 403001, Message: "subscription access expired", HttpCode: http.StatusForbidden}
	ErrAccessLimitExceeded       = customerrors.ExternalError{Code: 403002, Message: "limit of access exceeded", HttpCode: http.StatusForbidden}
	ErrAccessNotExist            = customerrors.ExternalError{Code: 403003, Message: "access key does not exist", HttpCode: http.StatusForbidden}
)
