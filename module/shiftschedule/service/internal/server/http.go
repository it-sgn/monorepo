package server

import (
	// v1 "mall-go/api/shiftschedule/service/v1"
	v1 "mall-go/api/shiftschedule/service/v1"
	"mall-go/module/shiftschedule/service/internal/service"

	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"

	"mall-go/module/shiftschedule/service/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	// _ "github.com/go-kratos/swagger-ui/swag"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, as *service.ShiftScheduleService, logger log.Logger) *http.Server {
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

	v1.RegisterShiftScheduleHTTPServer(srv, as)
	// v1.register

	// empv1.RegisterEmployersHTTPServer(srv, as)
	// bioV1.RegisterBiometricHTTPServer(srv, bs)

	return srv
}
