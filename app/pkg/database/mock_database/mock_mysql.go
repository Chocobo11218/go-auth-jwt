package mock_database

import (
	"errors"

	"github.com/DATA-DOG/go-sqlmock"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	conn, sqlMock, _ := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, nil, errors.Join(err, errors.New("can't initialize entity session"))
	}

	return db, sqlMock, nil
}