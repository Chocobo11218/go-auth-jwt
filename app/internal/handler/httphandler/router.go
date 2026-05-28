package httphandler

import (
	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	config             *config.AppConfig
	server             *echo.Echo
	healthCheckHandler IHealthCheckHandler
	// SaleSummaryReportHandler Handler
}

func NewHttpServer(
	config *config.AppConfig,
	server *echo.Echo,
	healthHealthCheckHandler IHealthCheckHandler,
	// saleSummaryReportHandler Handler,
) *HttpServer {
	httpServer := &HttpServer{
		config:             config,
		server:             server,
		healthCheckHandler: healthHealthCheckHandler,
		// SaleSummaryReportHandler: saleSummaryReportHandler,
	}

	httpServer.InitRoutes()

	return httpServer
}

func (s *HttpServer) InitRoutes() {
	e := s.server

	e.GET("/health", s.healthCheckHandler.HealthCheck)
	e.GET("/ready", s.healthCheckHandler.ReadinessCheck)

	rootCtx := e.Group("")

	api := rootCtx.Group("/api")

	v1 := api.Group("/v1")
	_ = v1
	// v1.POST("/list", s.SaleSummaryReportHandler.GetReportList)
}

func (s *HttpServer) Start(address string) error {
	return s.server.Start(address)
}

func (s *HttpServer) Server() *echo.Echo {
	return s.server
}
