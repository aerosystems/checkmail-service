package HTTPServer

type CheckHandler struct {
	*BaseHandler
	inspectUsecase InspectUsecase
}

func NewCheckHandler(
	baseHandler *BaseHandler,
	inspectUsecase InspectUsecase,
) *CheckHandler {
	return &CheckHandler{
		BaseHandler:    baseHandler,
		inspectUsecase: inspectUsecase,
	}
}
