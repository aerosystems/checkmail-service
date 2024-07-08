package RpcServer

type InspectRPCPayload struct {
	Domain   string
	ClientIp string
}

func (s Server) Inspect(payload InspectRPCPayload, resp *string) error {
	domainType, err := s.inspectUsecase.InspectData(payload.Domain, payload.ClientIp, "")
	if err != nil {
		return err
	}

	*resp = domainType.String()
	return nil
}
