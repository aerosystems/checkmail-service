package rpcServer

type InspectRPCPayload struct {
	Domain   string
	ClientIp string
}

func (cs RPCServer) Inspect(payload InspectRPCPayload, resp *string) error {
	domainType, err := cs.inspectUsecase.InspectData(payload.Domain, payload.ClientIp, "")
	if err != nil {
		return err.Error()
	}

	*resp = *domainType
	return nil
}
