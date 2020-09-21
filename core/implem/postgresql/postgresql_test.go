package postgresql_test

import (
	"home-broker/core/implem/postgresql"
	postgresqltests "home-broker/tests/postgresql"
	"reflect"
	"testing"
)

func TestNewDB(t *testing.T) {
	expected := postgresqltests.GetDB()
	db := postgresql.NewDB(expected.Host, expected.Port, expected.User, expected.Password, expected.DBName)
	if !reflect.DeepEqual(db, expected) {
		t.Errorf("received %v, expected %v", db, &expected)
	}
}
