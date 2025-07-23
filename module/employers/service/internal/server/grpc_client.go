package server

import (
	"context"

	biometricV1 "mall-go/api/biometrics/service/v1"
	departmentV1 "mall-go/api/department/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	grpcx "github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewBiometricClient menggunakan registry.Discovery langsung
func NewBiometricClient(discovery registry.Discovery) biometricV1.BiometricClient {
	conn, err := grpcx.DialInsecure(
		context.Background(),
		grpcx.WithEndpoint("discovery:///biometric.service"),
		grpcx.WithDiscovery(discovery),
	)
	if err != nil {
		panic(err)
	}
	return biometricV1.NewBiometricClient(conn)
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
