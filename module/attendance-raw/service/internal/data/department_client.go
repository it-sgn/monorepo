package data

import (
	"context"
	v1 "mall-go/api/department/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewDepartmentClient(r registry.Discovery) (v1.DepartmentClient, error) {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///department.service"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		return nil, err
	}
	return v1.NewDepartmentClient(conn), nil
}
