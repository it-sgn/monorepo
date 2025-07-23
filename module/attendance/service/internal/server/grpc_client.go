package server

import (
	"context"

	// biometricV1 "mall-go/api/biometric/service/v1"
	departmentV1 "mall-go/api/department/service/v1"
	employersV1 "mall-go/api/employers/service/v1"

	"github.com/go-kratos/kratos/v2/registry"
	grpcx "github.com/go-kratos/kratos/v2/transport/grpc"
)

// // NewBiometricClient menggunakan registry.Discovery langsung
// func NewBiometricClient(discovery registry.Discovery) biometricV1.BiometricClient {
// 	conn, err := grpcx.DialInsecure(
// 		context.Background(),
// 		grpcx.WithEndpoint("discovery:///biometric.service"),
// 		grpcx.WithDiscovery(discovery),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return biometricV1.NewBiometricClient(conn)
// }

// NewBiometricClient menggunakan registry.Discovery langsung
func NewEmployersClient(discovery registry.Discovery) employersV1.EmployersClient {
	conn, err := grpcx.DialInsecure(
		context.Background(),
		grpcx.WithEndpoint("discovery:///employers.service"),
		grpcx.WithDiscovery(discovery),
	)
	if err != nil {
		panic(err)
	}
	return employersV1.NewEmployersClient(conn)
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
