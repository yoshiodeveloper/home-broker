package postgresql

import (
	"fmt"
	"home-broker/core"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB handles database connections for PostgreSQL.
type DB struct {
	core.DBInterface
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	gormDB   *gorm.DB
}

// NewDB creates a new PostgreSQL database client.
func NewDB(host string, port int, user string, password string, dbName string) DB {
	return DB{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  "disable",
	}
}

// Open connects to the database.
func (db *DB) Open() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.SetDB(gormDB)
	return err
}

// SetDB sets a pointer to a database connection.
func (db *DB) SetDB(gormDB *gorm.DB) {
	db.gormDB = gormDB
}

// GetDB returns a pointer to a GORM DB connection.
func (db DB) GetDB() *gorm.DB {
	return db.gormDB
}
