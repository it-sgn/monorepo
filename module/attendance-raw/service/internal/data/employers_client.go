package data

import (
	"context"
	v1 "mall-go/api/employers/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewEmployersClient(r registry.Discovery) (v1.EmployersClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///employers.service"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		return nil, err
	}
	return v1.NewEmployersClient(conn), nil
}
