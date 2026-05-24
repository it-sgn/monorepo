package server

import (
	"context"

	departmentV1 "mall-go/api/department/service/v1"
	employerV1 "mall-go/api/employers/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	grpcx "github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewBiometricClient menggunakan registry.Discovery langsung
func NewEmployersClient(discovery registry.Discovery) employerV1.EmployersClient {
	conn, err := grpcx.DialInsecure(
		context.Background(),
		grpcx.WithEndpoint("discovery:///employers.service"),
		grpcx.WithDiscovery(discovery),
	)
	if err != nil {
		panic(err)
	}
	return employerV1.NewEmployersClient(conn)
}

func NewDepartmentClient(discovery registry.Discovery) departmentV1.DepartmentClient {
	conn, err := grpcx.DialInsecure(
		context.Background(),
		grpcx.WithEndpoint("discovery:///department.service"),
		grpcx.WithDiscovery(discovery),
	)
	if err != nil {
		panic(err)
	}
	return departmentV1.NewDepartmentClient(conn)
}
