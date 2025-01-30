package adapters

type ExternalApiError struct {
	Code     int
	Message  string
	HttpCode int
}

func (e ExternalApiError) Error() string {
	return e.Message
}
