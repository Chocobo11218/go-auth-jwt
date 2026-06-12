package httphandler

import (
	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/middleware"
	"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
)

type HttpServer struct {
	config             *config.AppConfig
	server             *echo.Echo
	healthCheckHandler IHealthCheckHandler
	authHandler        AuthHandler
	userHandler        UserHandler
}

func NewHttpServer(
	config *config.AppConfig,
	server *echo.Echo,
	healthHealthCheckHandler IHealthCheckHandler,
	authHandler AuthHandler,
	userHandler UserHandler,
) *HttpServer {
	httpServer := &HttpServer{
		config:             config,
		server:             server,
		healthCheckHandler: healthHealthCheckHandler,
		authHandler:        authHandler,
		userHandler:        userHandler,
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

	v1.POST("/register", s.authHandler.Register)

	// rate-limit login: 5 requests/min per IP
	//loginLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5))
	v1.POST("/login", s.authHandler.Login)

	jwtAuth := middleware.JWTAuth(s.config)
	v1.GET("/me", s.userHandler.GetMe, jwtAuth)
}

func (s *HttpServer) Start(address string) error {
	return s.server.Start(address)
}

func (s *HttpServer) Server() *echo.Echo {
	return s.server
}
