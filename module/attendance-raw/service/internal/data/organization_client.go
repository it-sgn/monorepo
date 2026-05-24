package data

import (
	"context"
	v1 "mall-go/api/organization/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewOrganizationClient(r registry.Discovery) (v1.OrganizationServiceClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///organization.service"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		return nil, err
	}
	return v1.NewOrganizationServiceClient(conn), nil
}
