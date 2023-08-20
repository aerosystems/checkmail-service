package RPCServer

type RPCInspectPayload struct {
	Domain   string
	ClientIp string
}

func (cs *CheckmailServer) Inspect(payload RPCInspectPayload, resp *string) error {
	domainType, err := cs.InspectService.InspectData(payload.Domain, payload.ClientIp)
	if err != nil {
		return err.Error()
	}

	*resp = *domainType
	return nil
}
