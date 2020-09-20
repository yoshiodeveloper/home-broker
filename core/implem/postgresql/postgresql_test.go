package postgresql_test

import (
	"home-broker/core/implem/postgresql"
	testpostgresql "home-broker/test/postgresql"
	"reflect"
	"testing"
)

func TestNewDB(t *testing.T) {
	expected := testpostgresql.GetDB()
	db := postgresql.NewDB(expected.Host, expected.Port, expected.User, expected.Password, expected.DBName)
	if !reflect.DeepEqual(db, expected) {
		t.Errorf("received %v, expected %v", db, &expected)
	}
}
