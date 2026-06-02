package mysql_database

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"os"
	"time"
)

type MysqlConfig struct {
	Name                string
	SSLMode             string
	MaxOpenConns        *int
	MaxIdleConns        *int
	ConnMaxLifetimeHour time.Duration
	MysqlHost           string
	MysqlPort           string
	MysqlUser           string
	MysqlPassword       string
	Loc                 *time.Location
}

func Connect(config *MysqlConfig) (*gorm.DB, error) {
	if config.MysqlHost == "" {
		return nil, errors.New("host is required")
	}
	if config.MysqlPort == "" {
		return nil, errors.New("port is required")
	}
	if config.MysqlUser == "" {
		return nil, errors.New("user is required")
	}
	if config.MysqlPassword == "" {
		return nil, errors.New("password is required")
	}
	if config.Name == "" {
		return nil, errors.New("entity name is required")
	}
	if config.SSLMode == "" {
		config.SSLMode = "true"
	}

	if config.MaxOpenConns == nil {
		config.MaxOpenConns = new(int)
		*config.MaxOpenConns = 20
	}
	if config.MaxIdleConns == nil {
		config.MaxIdleConns = new(int)
		*config.MaxIdleConns = 10
	}

	var loc string
	if config.Loc != nil {
		loc = url.QueryEscape(config.Loc.String())
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&charset=utf8&loc=%s",
		config.MysqlUser,
		config.MysqlPassword,
		config.MysqlHost,
		config.MysqlPort,
		config.Name,
		loc,
	)

	newLogrus := logrus.New()
	newLogrus.SetFormatter(&logrus.JSONFormatter{})
	newLogrus.SetOutput(os.Stdout)

	gDB, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, errors.Join(err, errors.New("can't initialize entity session"))
	}

	mySQLDb, err := gDB.DB()
	if err != nil {
		return nil, errors.Join(err, errors.New("can't get mysql"))
	}

	mySQLDb.SetMaxOpenConns(*config.MaxOpenConns)
	mySQLDb.SetMaxIdleConns(*config.MaxIdleConns)
	mySQLDb.SetConnMaxLifetime(time.Hour * config.ConnMaxLifetimeHour)

	if err := mySQLDb.Ping(); err != nil {
		return nil, errors.Join(err, errors.New("mysql ping error"))
	}

	newLogrus.Infof("connect database %s:%s/%s successfully", config.MysqlHost, config.MysqlPort, config.Name)

	return gDB, nil
}
