package adapters

import (
	"context"
	"crypto/tls"
	"github.com/aerosystems/checkmail-service/internal/entities"
	"github.com/aerosystems/common-service/gen/protobuf/lookup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type LookupAdapter struct {
	client lookup.LookupServiceClient
}

func NewLookupAdapter(address string) (*LookupAdapter, error) {
	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30,
			Timeout: 30,
		}),
	}
	if address[len(address)-4:] == ":443" {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, err
	}
	return &LookupAdapter{
		client: lookup.NewLookupServiceClient(conn),
	}, nil
}

func (la LookupAdapter) Lookup(ctx context.Context, domain string) (entities.Type, error) {
	resp, err := la.client.Lookup(ctx, &lookup.LookupRequest{Domain: domain})
	if err != nil {
		return entities.UndefinedType, err
	}
	return entities.DomainTypeFromString(resp.DomainType), nil
}
