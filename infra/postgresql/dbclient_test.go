package postgresql_test

import (
	"home-broker/infra/postgresql"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetTestDBClient() postgresql.DBClient {
	dbClient := postgresql.NewDBClient(
		"Host",
		1234,
		"User",
		"Password",
		"DBName",
		"SSLMode",
	)
	return *dbClient
}

func GetMockedDBClient() (postgresql.DBClient, sqlmock.Sqlmock) {
	dbClient := GetTestDBClient()

	mockedDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockedDB,
	}), &gorm.Config{Logger: newLogger})

	dbClient.SetDB(gormDB)
	return dbClient, mock
}

func TestNewDBClient(t *testing.T) {
	expected := GetTestDBClient()
	dbClient := postgresql.NewDBClient(expected.Host, expected.Port, expected.User, expected.Password, expected.DBName, expected.SSLMode)
	if !reflect.DeepEqual(dbClient, &expected) {
		t.Errorf("received %v, expected %v", dbClient, &expected)
	}
}
