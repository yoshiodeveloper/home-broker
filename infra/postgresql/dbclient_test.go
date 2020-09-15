package postgresql_test

import (
	"home-broker/infra/postgresql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: mockedDB,
	}), &gorm.Config{})

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
