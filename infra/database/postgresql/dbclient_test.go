package postgresql

import (
	"reflect"
	"testing"
)

const (
	failTestMSG = "received %v, expected %v"
)

func TestNewDBClient(t *testing.T) {
	expected := DBClient{
		Host:     "Host",
		Port:     1234,
		User:     "User",
		Password: "Password",
		DBName:   "DBName",
		SSLMode:  "SSLMode",
	}
	dbClient := NewDBClient(expected.Host, expected.Port, expected.User, expected.Password, expected.DBName, expected.SSLMode)
	if !reflect.DeepEqual(dbClient, expected) {
		t.Errorf(failTestMSG, dbClient, expected)
	}
}
