package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/handler/httphandler"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/repository"
	"github.com/Chocobo11218/go-auth-jwt/app/internal/service"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/configurer"
	mysql_database "github.com/Chocobo11218/go-auth-jwt/app/pkg/database"

	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

func main() {
	// create root context
	ctx, cancel := context.WithCancel(context.Background())

	// load config
	conf := configurer.LoadConfig()

	// load secret env
	configurer.LoadSecret(&conf.Secret)
	fmt.Printf("%+v", conf)

	// init logger
	logger.Initialize()

	// validate config
	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		logger.Error(ctx, "fail validate loaded config", zap.Error(err))
	}

	if err := config.SetTimeZone(conf.Server.TimeZone); err != nil {
		logger.Error(ctx, "fail set timezone", zap.Error(err))
	}

	// database
	db, err := mysql_database.Connect(&mysql_database.MysqlConfig{
		Name:                conf.DB.Name,
		SSLMode:             conf.DB.SSLMode,
		MaxOpenConns:        &conf.DB.MaxOpenConnection,
		MaxIdleConns:        &conf.DB.MaxIdleConnection,
		ConnMaxLifetimeHour: conf.DB.ConnectionMaxLifetimeHour,
		MysqlHost:           conf.Secret.DBHost,
		MysqlPort:           conf.Secret.DBPort,
		MysqlUser:           conf.Secret.DBUser,
		MysqlPassword:       conf.Secret.DBPassword,
		Loc:                 time.Local,
	})
	fmt.Println("time: ", time.Now())

	if err != nil {
		logger.Error(ctx, "fail connect database", zap.Error(err))
		os.Exit(1)
	}
	defer func() {
		if db == nil {
			return
		}
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}()

	fmt.Println("db:", db)

	// init repository
	userRepo := repository.NewUserRepository(db)

	// init service
	jwtExpire, err := time.ParseDuration(conf.JWT.Expire)
	if err != nil {
		logger.Error(ctx, "fail parse jwt expire duration", zap.Error(err))
		os.Exit(1)
	}

	authService := service.NewAuthService(
		userRepo,
		conf.Secret.JWTSecret,
		jwtExpire,
	)
	//authService := service.NewAuthServiceMock()

	// init handler
	authHandler := httphandler.NewAuthHandler(authService)

	// init http server
	httpServer := httphandler.NewHttpServer(
		conf,
		echo.New(),
		httphandler.NewHealthCheckHandler(),
		authHandler,
	)

	go func() {
		if err := httpServer.Start(conf.Server.Address); err != nil {
			panic(errors.Wrap(err, "http server start error"))
		}
	}()

	monitorGraceful(ctx, cancel, httpServer.Server())
}

func monitorGraceful(
	ctx context.Context,
	cancel context.CancelFunc,
	httpServer *echo.Echo,
) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	select {
	case <-ctx.Done():
		logger.Info(ctx, "MonitorGraceful",
			zap.String("type", "server"),
			zap.String("msg", "MonitorGraceful - Terminating: context cancelled"),
		)
	case s := <-sigterm:
		logger.Info(ctx, "MonitorGraceful",
			zap.String("type", "server"),
			zap.String("msg", fmt.Sprintf("MonitorGraceful - Terminating: via signal %v", s)),
		)
	}

	cancel()
	newCtx, newCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer newCancel()

	if httpServer != nil {
		if err := httpServer.Shutdown(newCtx); err != nil {
			logger.Error(newCtx, "MonitorGraceful",
				zap.String("type", "server"),
				zap.String("msg", fmt.Sprintf("MonitorGraceful - Terminating: shutdown http server error %v", err)),
			)
		} else {
			logger.Info(newCtx, "MonitorGraceful",
				zap.String("type", "server"),
				zap.String("msg", "MonitorGraceful - Terminating: shutdown http server success"),
			)
		}
	}
}
