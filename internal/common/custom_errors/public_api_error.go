package CustomErrors

type PublicApiError struct {
	Code     int
	Message  string
	HttpCode int
}

func (e PublicApiError) Error() string {
	return e.Message
}

var (
	ErrEmailNotValid      = PublicApiError{400001, "email address does not valid", 400}
	ErrDomainNotValid     = PublicApiError{400002, "domain does not valid", 400}
	ErrDomainNotExist     = PublicApiError{400003, "domain does not exist", 400}
	ErrDomainTrustedTypes = PublicApiError{400003, "domain type does not exist in trusted types", 400}
	ErrDomainCoverage     = PublicApiError{400004, "domain coverage does not exist in trusted coverages", 400}
)
