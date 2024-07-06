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
	ErrEmailNotValid  = PublicApiError{400001, "email address does not valid", 400}
	ErrDomainNotExist = PublicApiError{400003, "domain does not exist", 400}
)
