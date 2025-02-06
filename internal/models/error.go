package models

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
)

var (
	ErrEmailNotValid      = customerrors.ExternalError{Code: 400001, Message: "email address does not valid", HttpCode: 400}
	ErrDomainNotValid     = customerrors.ExternalError{Code: 400002, Message: "domain does not valid", HttpCode: 400}
	ErrDomainNotExist     = customerrors.ExternalError{Code: 400003, Message: "domain does not exist", HttpCode: 400}
	ErrDomainTrustedTypes = customerrors.ExternalError{Code: 400003, Message: "domain type does not exist in trusted types", HttpCode: 400}
	ErrDomainCoverage     = customerrors.ExternalError{Code: 400004, Message: "domain coverage does not exist in trusted coverages", HttpCode: 400}
)
