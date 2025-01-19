package HTTPServer

type AccessHandler struct {
	accessUsecase AccessUsecase
}

func NewAccessHandler(accessUsecase AccessUsecase) *AccessHandler {
	return &AccessHandler{accessUsecase: accessUsecase}
}
