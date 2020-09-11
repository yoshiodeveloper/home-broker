package mock

import (
	"broker-dealer/infra"
)

// DBClient handles database connections.
type DBClient struct {
	infra.DBClient
	db *Database
}

// NewDBClient creates a new database client.
func NewDBClient() *DBClient {
	db := NewDatabase()
	return &DBClient{db: &db}
}

// SetDB sets a pointer to a database connection.
func (dbClient *DBClient) SetDB(db *Database) {
	dbClient.db = db
}

// Open connects to the database.
func (dbClient *DBClient) Open() error {
	return nil
}

// GetDB returns a pointer to a GORM DB connection.
func (dbClient *DBClient) GetDB() *Database {
	return dbClient.db
}
