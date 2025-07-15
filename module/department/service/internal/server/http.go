package server

import (
	empv1 "mall-go/api/department/service/v1"
	"mall-go/module/department/service/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	"mall-go/module/department/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	// _ "github.com/go-kratos/swagger-ui/swag"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, ds *service.DepartmentService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.Validator(),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	srv.HandlePrefix("/q/", openapiv2.NewHandler())

	// v1.RegisterEmployersHTTPServer(srv, as)

	empv1.RegisterDepartmentHTTPServer(srv, ds)
	// bioV1.RegisterBiometricHTTPServer(srv, bs)

	return srv
}
